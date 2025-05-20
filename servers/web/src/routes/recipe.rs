use std::env;

use askama::Template;
use axum::{extract::Path, http::request, response::{Html, IntoResponse}};
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use tracing::info;

use crate::{utils::{self, objects::Recipe}, AppError};

#[derive(Template)]
#[template(path = "pages/recipe.html")]
struct RecipePage {
    recipe: Recipe,
}
#[derive(Serialize)]
struct Request {
    #[serde(alias="URL")]
    url: String
}

pub(crate) async fn recipe_handler(Path(url): Path<String>) -> Result<impl IntoResponse, AppError> {
    info!("Getting recipe for {}", url);

    let mut json = utils::extract::extract_json_from_url(&url).await?;

    // TODO: handle more than one recipe
    if json.is_empty() {
        return Err(AppError::MiscError(String::from("No json schema found")))
    }
    let web_recipe_json = json.remove(0);

    let recipe_json: Recipe = web_recipe_json.into();

    // let request = Request { url };

    // let data_port = env::var("DATA_PORT").expect("No data port provided");
    // let client = reqwest::Client::new();
    // let response: Value = client.post(format!("http://data_service:{}/import-url", data_port))
    //     .json(&request)
    //     .send()
    //     .await?
    //     .json()
    //     .await?;

    // info!("Recieved response from datea server: {}", response);

    // let recipe_page: RecipePage = serde_json::from_value(response)?;
    let recipe_page = RecipePage {
        recipe: recipe_json,
    };

    Ok(Html(recipe_page.render().expect("recipe should render")))

}