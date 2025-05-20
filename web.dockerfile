# Build the rust project, it uses a Debain OS
FROM rust:1.87.0 AS builder

WORKDIR /

RUN mkdir -p ./servers/web/src

# Cache dependencies work around, copying over all files triggers a complete rebuuild as cargo.toml is new.
# Copying individually allows us only rebuild dependencies if they have changed
COPY Cargo.toml .
COPY Cargo.lock .
COPY servers/web/Cargo.toml /servers/web
COPY servers/web/Cargo.lock /servers/web

# Create a dummy file
RUN echo 'fn main() { println!("hello"); }' > ./servers/web/src/main.rs
RUN cargo build --release -p openchef-web
RUN rm -rf ./servers/web/src

COPY servers/web/src/ /servers/web/src/
COPY servers/web/templates/ /servers/web/templates/

# Force update modification times to make cargo rebuild them
RUN touch -a -m ./servers/web/src/main.rs

RUN cargo build --release -p openchef-web

# Deploy the application
FROM ubuntu:25.10

WORKDIR /
RUN apt-get update && apt-get install -y ca-certificates

COPY /servers/web/static /static
COPY --from=builder /target/release/openchef-web .

ENTRYPOINT [ "./openchef-web" ]