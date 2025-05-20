use std::env;

use askama::Template;
use axum::{http::StatusCode, response::{Html, IntoResponse, Response}, routing::get, Router};
use routes::recipe::recipe_handler;
use thiserror::Error;
use tower_http::services::ServeDir;
use tracing::{info, Level};
use tracing_subscriber::FmtSubscriber;

mod routes;
mod utils;

#[derive(Template)]
#[template(path = "pages/home.html")]
struct Home {}

async fn home() -> impl IntoResponse {
    let home = Home {};

    Html(home.render().expect("Why isnt the home page rendering?!"))
}

async fn ping() -> impl IntoResponse {
    "Pong"
}

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Error making outgoing request {0}")]
    OutgoingRequestError(#[from] reqwest::Error),
    #[error("Error creating regex {0}")]
    RegexError(#[from] regex::Error),
    #[error("Error [de?]serialising data {0}")]
    JsonError(#[from] serde_json::Error),
    #[error("{0}")]
    MiscError(String)
}

impl From<String> for AppError {
    fn from(value: String) -> Self {
        Self::MiscError(value)
    }
}

// TODO: implement a error template
impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("Something went wrong: {:?}", self),
        )
            .into_response()
    }
}


#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let subscriber = FmtSubscriber::builder()
        .with_max_level(Level::DEBUG)
        .finish();
    tracing::subscriber::set_global_default(subscriber).expect("setting default subscriber failed");

    info!("Starting the web server...");

    let router = Router::new()
        .route("/", get(home))
        .route("/ping", get(ping))
        .route("/recipe/{*url}", get(recipe_handler))
        .nest_service("/static", ServeDir::new("./static"));

    let port = env::var("WEB_PORT")?;

    let http_addr: std::net::SocketAddr = format!("0.0.0.0:{}", port).parse()?;

    info!("Server is running on {}", http_addr);

    let http_server = axum_server::bind(http_addr).serve(router.into_make_service());

    http_server.await?;

    Ok(())
}