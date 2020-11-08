

use clap::{ crate_description, crate_name, crate_version, App, Arg };


fn main() {
    let matches = App::new(crate_name!())
        .version(crate_version!())
        .about(crate_description!())
        .arg(
            Arg::with_name("client")
                .short("c")
                .long("client")
                .value_name("HOSTNAME")
                .help("Run in client mode with server at HOSTNAME")
                .takes_value(true)
                .required_unless("server")
                .conflicts_with("server"),
        )
        .arg(
            Arg::with_name("server")
                .short("s")
                .long("server")
                .help("Run in server mode")
                .takes_value(false)
                .required_unless("client")
                .conflicts_with("client"),
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
        println!("Running in client mode");
        //handle_client(matches.value_of("client").unwrap_or(""));
    } else if matches.is_present("server") {
        // handle_server(
        //     matches
        //         .value_of("port")
        //         .unwrap_or(&format!("{}", constants::DEFAULT_PORT)),
        // );
        println!("Running in server mode");
    }
}