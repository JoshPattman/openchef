use std::env;

use askama::Template;
use axum::{extract::Path, http::request, response::{Html, IntoResponse}};
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use tracing::info;

use crate::{get_db_connection, utils::{self, db::{add_parser_error, ParserError}, objects::Recipe}, AppError};

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

    let mut recipes = utils::extract::extract_json_from_url(&url).await?;

    // TODO: handle more than one recipe
    let recipe = recipes.remove(0);

    let recipe_page = RecipePage {
        recipe,
    };

    Ok(Html(recipe_page.render().expect("recipe should render")))

}