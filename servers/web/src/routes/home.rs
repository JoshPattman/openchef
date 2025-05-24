use askama::Template;
use axum::{response::{Html, IntoResponse, Redirect}, Form};
use serde::Deserialize;

use crate::AppError;


#[derive(Template)]
#[template(path = "pages/home.html")]
struct Home {}

pub(crate) async fn home() -> impl IntoResponse {
    let home = Home {};

    Html(home.render().expect("Why isnt the home page rendering?!"))
}

#[derive(Deserialize)]
pub(crate) struct RecipeForm {
    recipe_url: String,
}

pub(crate) async fn navigate(Form(form): Form<RecipeForm>) -> Result<impl IntoResponse, AppError> {
    let url = form.recipe_url.trim();
    if url.is_empty() {
        return Err(String::from("No Url Given").into())
    }
    
    Ok(([("HX-Redirect", format!("/recipe/{}", url))], "OK"))   
}