mod book;

fn main() {
    let book = book::get_book();

    for x in &book {
        println!("{}", &x);
    }
}
