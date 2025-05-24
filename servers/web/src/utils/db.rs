use std::{env, str::FromStr, sync::OnceLock};

use sqlx::{prelude::FromRow, sqlite::{SqliteConnectOptions, SqlitePool, SqlitePoolOptions}, Executor, Sqlite};
use tracing::info;

use crate::AppError;

pub(crate) async fn get_db_connection() -> Result<SqlitePool, AppError> {
    let db_url = env::var("DATA_PERSIST_PATH")?;
    let pool = SqlitePoolOptions::new()
        .max_connections(3)
        .connect(&format!("sqlite://{}/web-server.db", db_url))
        .await?;
    
    Ok(pool)
}

pub(crate) async fn create_tables(pool: &SqlitePool) -> Result<(), AppError> {
    sqlx::query(
        r#"
        CREATE TABLE IF NOT EXISTS parser_errors (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            url TEXT NOT NULL,
            error TEXT NOT NULL
        );
        "#
    )
        .execute(pool)
        .await?;

    Ok(())
}

#[derive(sqlx::FromRow, Debug)]
pub(crate) struct ParserError {
    pub url: String,
    pub error: String,
}

impl ParserError {
    pub fn new(url: String, error: String) -> Self {
        Self { url, error }
    }
}

pub(crate) async fn add_parser_error(pool: &SqlitePool, error: ParserError) -> Result<(), AppError> {
    // Check if the entry already exists
    let existing: i64 = sqlx::query_scalar(
        r#"
        SELECT COUNT(*)
        FROM parser_errors
        WHERE url = ? AND error = ?;
        "#
    )
    .bind(&error.url)
    .bind(&error.error)
    .fetch_one(pool)
    .await?;

    if existing > 0 {
        info!("Parser error already in database {:?}", error);
        return Ok(()); // Entry already exists, do nothing
    }

    // Insert the new entry
    sqlx::query(
        r#"
        INSERT INTO parser_errors (url, error)
        VALUES (?, ?);
        "#
    )
    .bind(&error.url)
    .bind(&error.error)
    .execute(pool)
    .await?;

    Ok(())
}

pub(crate) async fn get_parser_errors(pool: &SqlitePool) -> Result<Vec<ParserError>, AppError> {
    let errors = sqlx::query_as::<_, ParserError>(
        r#"
        SELECT url, error
        FROM parser_errors;
        "#
    )
    .fetch_all(pool)
    .await?;

    Ok(errors)
}

pub(crate) async fn remove_url_from_parser_errors(pool: &SqlitePool, url: &str) -> Result<(), AppError> {
    sqlx::query(
        r#"
        DELETE FROM parser_errors
        WHERE url = ?;
        "#
    )
    .bind(url)
    .execute(pool)
    .await?;

    Ok(())
}


