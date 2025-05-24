use std::env;

use axum::{http::StatusCode, response::{IntoResponse, Response}, routing::{get, post}, Router};
use routes::{home::{home, navigate}, recipe::recipe_handler};
use thiserror::Error;
use tower_http::services::ServeDir;
use tracing::{info, Level};
use tracing_subscriber::FmtSubscriber;
use utils::db::{create_tables, get_db_connection};

mod routes;
mod utils;

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
    #[error("Cannot find var {0}")]
    VarError(#[from] env::VarError),
    #[error("SQL error {0}")]
    SqlError(#[from] sqlx::Error),
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

    let pool = get_db_connection().await?;
    create_tables(&pool).await?;

    info!("Starting the web server...");

    let router = Router::new()
        .route("/", get(home))
        .route("/navigate", post(navigate))
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