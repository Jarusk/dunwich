pub struct DunwichServer {
    port: u16,
}

impl DunwichServer {
    pub fn new(port: u16) -> DunwichServer {
        DunwichServer { port }
    }

    pub fn run(&self) {
        println!("{:?}", self.port);
    }
}
