use std::fmt::Display;

use serde::{Deserialize, Serialize};
use tracing::error;

use super::extract::{WebImage, WebRecipe, WebRecipeInstructions, WebRecipeYield};


#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub(crate) struct Recipe {
    pub name: String,
    pub description: String,
    pub keywords: String,
    pub image: Image,
    pub prep_time: String,
    pub cook_time: String,
    pub ingredients: Vec<Ingredient>,
    pub instructions: Vec<String>,
    pub yields: Yield,
}

impl From<WebRecipe> for Recipe {
    fn from(value: WebRecipe) -> Self {
        Recipe { 
            name: value.name,
            description: value.description,
            keywords: value.keywords,
            image: value.image.into(),
            prep_time: value.prep_time,
            cook_time: value.cook_time,
            ingredients: value.recipe_ingredient.into_iter().map(|x| x.into()).collect(),
            instructions: value.recipe_instructions.into(),
            yields: value.recipe_yield.into(),
        }
    }
}

impl From<WebRecipeInstructions> for Vec<String> {
    fn from(value: WebRecipeInstructions) -> Self {
        match value {
            WebRecipeInstructions::String(s) => s.split(',').map(|s| s.trim().to_string()).collect(),
            WebRecipeInstructions::StringArray(arr) => arr,
            WebRecipeInstructions::ObjectArray(objects) => objects.into_iter().map(|obj| obj.text).collect(),
        }
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct Yield {
    pub min: usize,
    pub max: usize,
}

impl Display for Yield {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        if self.min == self.max {
            write!(f, "{}", self.min)
        } else {
            write!(f, "{}-{}", self.min, self.max)
        }
    }
}

impl From<WebRecipeYield> for Yield {
    fn from(value: WebRecipeYield) -> Self {
        match value {
            WebRecipeYield::Number(y) => Yield { min: y, max: y },
            WebRecipeYield::String(s) => {
                match parse_yield_string(&s) {
                    Some(y) => y,
                    None => {
                        error!("Could not parse yield string {}", s);
                        Yield { min: 0, max: 0 }
                    }
                }
            },
            WebRecipeYield::Array(items) => {
                match items.iter().find_map(|s| parse_yield_string(s)) {
                    Some(y) => y,
                    None => {
                        error!("Could not any value in yield array {:?}", items);
                        Yield { min: 0, max: 0 }
                    }
                }
            },
        }
    }
}

fn parse_yield_string(s: &str) -> Option<Yield> {
    // TODO: find first number in string
    let num_pattern = r".*?(\d+)";
    let re = regex::Regex::new(num_pattern).expect("The regex is static, it should parse");

    re.captures(s).and_then(|c| {
        let v = c.get(1)?.as_str().parse().ok()?;
        Some(Yield { min: v, max: v })
    })
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct Image {
    pub url: String,
    pub height: usize,
    pub width: usize,
    pub caption: String,
}

impl From<WebImage> for Image {
    fn from(value: WebImage) -> Self {
        match value {
            WebImage::Object { width, height, url , caption} => Image {
                url,
                height,
                width,
                caption: caption.unwrap_or_default(),
            },
            WebImage::Array(urls) => Image {
                url: urls.get(0).cloned().unwrap_or_default(), // Use the first URL or default to an empty string
                height: 0,
                width: 0, 
                caption: String::new(),
            },
            WebImage::String(url) => Image {
                url,
                height: 0, 
                width: 0, 
                caption: String::new(), 
            },
        }
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct Ingredient {
    pub name: String,
    pub quantity: f64,
    pub metric: String,
}

impl From<String> for Ingredient {
    fn from(value: String) -> Self {
        // TODO: do this properly
        Ingredient { name: value, quantity: 1.0, metric: String::from("Unit") }
    }
}