extern crate clap;

use clap::{crate_authors, crate_description, crate_name, crate_version, App, Arg};
use std::net::SocketAddr;

mod book;
mod client;
mod server;

const DEFAULT_PORT: u16 = 5201;

fn main() {
    let matches = App::new(crate_name!())
        .version(crate_version!())
        .author(crate_authors!())
        .about(crate_description!())
        .arg(
            Arg::with_name("client")
                .short("c")
                .long("client")
                .value_name("HOSTNAME")
                .help("Run in client mode with server at HOSTNAME")
                .takes_value(true)
                .required_unless("server"),
        )
        .arg(
            Arg::with_name("server")
                .short("s")
                .long("server")
                .help("Run in server mode")
                .takes_value(false)
                .required_unless("client"),
        )
        .arg(
            Arg::with_name("port")
                .short("p")
                .long("port")
                .help("Set server port")
                .takes_value(true),
        )
        .arg(
            Arg::with_name("v")
                .short("v")
                .multiple(true)
                .help("Sets the level of verbosity"),
        )
        .get_matches();

    if matches.is_present("client") {
        handle_client(matches.value_of("client").unwrap_or(""));
    } else if matches.is_present("server") {
        handle_server(
            matches
                .value_of("port")
                .unwrap_or(&format!("{}", DEFAULT_PORT)),
        );
    }
}

fn handle_server(port: &str) {
    let port_parsed = match port.trim().parse::<u16>() {
        Ok(x) => x,
        Err(_e) => {
            eprintln!("Exiting: Invalid port specified ({})", &port);
            std::process::exit(1);
        }
    };

    server::DunwichServer::new(port_parsed).run();
}

fn handle_client(address: &str) {
    let mut tmp = address.trim().to_string();
    if !address.contains(':') {
        tmp = format!("{}:{}", tmp, DEFAULT_PORT);
    }
    let address_parsed = match tmp.parse::<SocketAddr>() {
        Ok(e) => e,
        Err(_e) => {
            eprintln!("Exiting: invalid ip:port address ({})", &address);
            std::process::exit(1);
        }
    };

    client::DunwichClient::new(address_parsed).run();
}
