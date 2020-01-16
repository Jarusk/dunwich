

struct MessageTypes {
    EndSession,
    Text
}

struct EndTestMessage { }


struct SessionMessage<'a> {
    length: u16,
    value: &'a str
}