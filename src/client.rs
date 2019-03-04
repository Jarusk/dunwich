use crate::book;

use std::io::prelude::*;
use std::net::{SocketAddr, TcpStream};
use std::str;
use std::time::Duration;

pub struct DunwichClient {
    address: SocketAddr,
    content: Vec<&'static str>
}

impl DunwichClient {
    pub fn new(address: SocketAddr) -> DunwichClient {
        DunwichClient { address: address , content: book::get_book()}
    }

    pub fn run(&self) {
        println!("Launching client connecting to {:?}", self.address);

        match TcpStream::connect_timeout(&self.address, Duration::new(5, 0)) {
            Ok(s) => {
                println!("Connected to server: {}", self.address);
                self.handle_connection(s);
            }
            Err(e) => {
                println!("Error in conneciton: {}", e);
            }
        }
    }

    fn handle_connection(&self, mut stream: TcpStream) {
        let mut tmp = [0u8;1024*1024];

        for x in 0..1000 {
            stream.read(&mut tmp).unwrap();
            //println!("{}", str::from_utf8(&tmp).unwrap());
            println!("{}", &x);
        }
    }
}