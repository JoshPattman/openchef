use std::env;

use askama::Template;
use axum::{extract::Path, http::request, response::{Html, IntoResponse}};
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use tracing::info;

use crate::AppError;

#[derive(Template, Deserialize)]
#[template(path = "pages/recipe.html")]
struct RecipePage {
    name: String,
    ingredients: Vec<Ingredient>,
    steps: Vec<String>,
    #[serde(alias="yield")]
    recipe_yield: usize,
}

#[derive(Deserialize)]
struct Ingredient {
    name: String,
    quantity: f64,
    metric: String,
}

#[derive(Serialize)]
struct Request {
    #[serde(alias="URL")]
    url: String
}

pub(crate) async fn recipe_handler(Path(url): Path<String>) -> Result<impl IntoResponse, AppError> {
    info!("Getting recipe for {}", url);

    let request = Request { url };

    let data_port = env::var("DATA_PORT").expect("No data port provided");
    let client = reqwest::Client::new();
    let response: Value = client.post(format!("http://data_service:{}/import-url", data_port))
        .json(&request)
        .send()
        .await?
        .json()
        .await?;

    info!("Recieved response from datea server: {}", response);

    let recipe_page: RecipePage = serde_json::from_value(response)?;

    Ok(Html(recipe_page.render().expect("recipe should render")))

}