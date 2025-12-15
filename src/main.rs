use std::env;

use toukei::{Config, FileReader};

fn main() {
    let args = env::args().collect::<Vec<String>>();

    let config = Config::build(&args).unwrap();

    let mut reader = FileReader::new(config);
    let _res = reader.run();

}
