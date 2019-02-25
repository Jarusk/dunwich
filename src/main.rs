extern crate clap;

use clap::{crate_authors, crate_description, crate_name, crate_version, App, Arg};

mod book;

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
            Arg::with_name("v")
                .short("v")
                .multiple(true)
                .help("Sets the level of verbosity"),
        )
        .get_matches();
}
