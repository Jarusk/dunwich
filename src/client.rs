use crate::constants;

use std::error::Error;
use std::io::prelude::*;
use std::net::{SocketAddr, TcpStream};
use std::sync::atomic::{AtomicUsize, Ordering};
use std::time::Duration;

static CLIENT_COUNT: AtomicUsize = AtomicUsize::new(0);

pub struct DunwichClient {
    id: usize,
    address: SocketAddr
}

impl DunwichClient {
    pub fn new(address: SocketAddr) -> DunwichClient {
        DunwichClient { id: CLIENT_COUNT.fetch_add(1, Ordering::Relaxed), address: address }
    }

    pub fn run(&self) {
        println!("Launching client {} connecting to {:?}", self.id, self.address);

        match TcpStream::connect_timeout(&self.address, Duration::new(constants::DEFAULT_CLIENT_TIMEOUT_SECONDS, 0)) {
            Ok(s) => {
                println!("Connected to server: {}", self.address);
                self.handle_connection(s);
                println!("")
            }
            Err(e) => {
                println!("Error in connection: {}", e.description());
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