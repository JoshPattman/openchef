use serde::{Deserialize, Serialize};
use serde_json::Value;
use tracing::{debug, error, info};

use crate::{utils::db::{add_parser_error, get_db_connection, remove_url_from_parser_errors, ParserError}, AppError};

use super::objects::Recipe;

// add fields as needed
#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub(crate) struct WebRecipe {
    pub name: String,
    pub description: String,
    pub keywords: Option<String>,
    pub image: WebImage,
    pub prep_time: String,
    pub cook_time: String,
    pub recipe_ingredient: Vec<String>,
    pub recipe_instructions: WebRecipeInstructions,
    pub recipe_yield: WebRecipeYield,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(untagged)]
pub(crate) enum WebImage {
    Object {
        width: usize,
        height: usize,
        #[serde(alias="contentUrl")]
        url: String,
        caption: Option<String>,
    },
    Array(Vec<String>),
    String(String),
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(untagged)]
pub(crate) enum WebRecipeInstructions {
    String(String),
    StringArray(Vec<String>),
    ObjectArray(Vec<WebRecipeInstructionsObject>)
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct WebRecipeInstructionsObject {
    pub text: String
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(untagged)]
pub(crate) enum WebRecipeYield {
    Number(usize),
    String(String),
    Array(Vec<String>),
}

pub(crate) async fn extract_json_from_url(url: &str) -> Result<Vec<Recipe>, AppError> {
    info!("Extracting json from {}", url);

    // Fetch the HTML content from the URL
    let response = reqwest::get(url).await?;
    let body = response.text().await?;

    info!("Found {} bytes in body", body.bytes().len());

    // Look for JSON-LD script tags
    let script_pattern = r#"(?s)<script[^>]*type="application\/ld\+json"[^>]*>(.*?)<\/script>"#;
    let re = regex::Regex::new(script_pattern)?;

    let mut recipe_jsons: Vec<Value> = vec![];

    let pool = get_db_connection().await?;

    // Extract the JSON content from all matching script tags
    for capture in re.captures_iter(&body).filter_map(|c| c.get(1)) {
        let content = capture.as_str();
        debug!("Found schema: {}", content);
        
        // if the match has no recipe type, we skip
        let recipe_pattern = r#""@type":\s*"Recipe""#;
        let recipe_re = regex::Regex::new(recipe_pattern)?;
        if recipe_re.find(content).is_none() {
            debug!("No recipe tag found in match");
            continue;
        }

        let parsed_json: Value = serde_json::from_str(content)?;

        // if object has a graph at top level, retrieve recipe from it
        if let Some(graph) = parsed_json.get("@graph").and_then(|g| g.as_array()) {
            debug!("Schema is graph type");

            for json in graph {
                if let Some(json_type) = json.get("@type").and_then(|t| t.as_str()) {
                    if json_type != "Recipe" {
                        continue;
                    }
                    recipe_jsons.push(json.clone());                        
                }
            }

            continue;
        }

        // top level recipe schema
        if "Recipe" == parsed_json.get("@type").and_then(|t| t.as_str()).unwrap_or("") {
            recipe_jsons.push(parsed_json);
            continue;
        }
        
        // if we get here, we should check the schema
        error!("New schema format found: {}", content);
        add_parser_error(&pool, ParserError::new(url.to_string(), "NEW SCHEMA FORMAT".to_string())).await?;
    }

    let mut web_recipes: Vec<WebRecipe> = vec![];
    for recipe_json in recipe_jsons.into_iter() {
        let web_recipe = match serde_json::from_value(recipe_json) {
            Ok(v) => v,
            Err(e) => {
                add_parser_error(&pool, ParserError::new(url.to_string(), format!("JSON PARSE ERROR: {}", e.to_string()))).await?;
                return Err(e.into());
            }
        };
        web_recipes.push(web_recipe);
    }

    if web_recipes.is_empty() {
        add_parser_error(&pool, ParserError::new(url.to_string(), "NO JSON SCHEMA FOUND".to_string())).await?;
        return Err(AppError::MiscError(String::from("No json schema found")))
    }

    
    // Recipes are valid from here onward
    let recipes = web_recipes.into_iter().map(|r| r.into()).collect();
    remove_url_from_parser_errors(&pool, url).await?;

    Ok(recipes)
}