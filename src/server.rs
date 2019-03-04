extern crate rand;

use std::io::prelude::*;
use std::net::{TcpListener, TcpStream};
use std::process;

use rand::Rng;

use crate::book;

pub struct DunwichServer {
    port: u16,
    content: Vec<&'static str>,
    content_len: usize
}

impl DunwichServer {
    pub fn new(port: u16) -> DunwichServer {
        let tmp = book::get_book();
        let tmp_len = tmp.len();
        DunwichServer { port: port, content: tmp, content_len: tmp_len }
    }

    pub fn run(&self) {
        println!("Launching server on port {:?}", self.port);
        let listener = match TcpListener::bind(format!("127.0.0.1:{}", self.port)) {
            Ok(l) => l,
            Err(e) => {
                eprintln!("Failed to setup server on port {}: {}", self.port, e);
                process::exit(1);
            }
        };

        for stream in listener.incoming() {
            match stream {
                Ok(stream) => {
                    println!("Accepted new client: {:?}", stream);
                    self.handle_client(stream);
                },
                Err(e) => {
                    eprintln!("Connection error: {}", e);
                }
            }
        }
    }

    fn handle_client(&self, mut stream: TcpStream) {

        let mut rng = rand::thread_rng();
        
        for x in 0..1000 {
            let i = rng.gen_range(0, self.content_len);
            stream.write(self.content[i].as_bytes()).unwrap();
            println!("{}", &x);
        }
    }
}