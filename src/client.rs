use std::net::SocketAddr;
use crate::book;

pub struct DunwichClient {
    address: SocketAddr,
}

impl DunwichClient {
    pub fn new(address: SocketAddr) -> DunwichClient {
        DunwichClient { address }
    }

    pub fn run(&self) {
        println!("{:?}", self.address);
    }
}
