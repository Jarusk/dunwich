use std::net::SocketAddr;

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
