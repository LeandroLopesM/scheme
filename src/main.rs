mod scheme;

use clap::Parser;
use log::{error, warn};
use reedline::{DefaultPrompt, DefaultPromptSegment, Reedline, Signal};

use crate::scheme::engine::Engine;

#[derive(Parser)]
#[command(about, version)]
struct Args {
    file: Option<String>,
}

fn main() {
    let args = Args::parse();
    let mut engine = Engine::new();

    colog::default_builder()
        .filter_level(if cfg!(debug_assertions) {
            log::LevelFilter::Trace
        } else {
            log::LevelFilter::Warn
        })
        .filter_module("rustyline", log::LevelFilter::Warn)
        .init();

    if let Some(_) = args.file {
        todo!("File interpretation")
    }

    let mut rl = Reedline::create();
    let mut prompt = DefaultPrompt::default();
    prompt.left_prompt = DefaultPromptSegment::Basic("λλλ ".to_string());
    prompt.right_prompt = DefaultPromptSegment::Empty;

    loop {
        match rl.read_line(&prompt) {
            Ok(Signal::Success(str)) => engine.run_str("REPL".to_string(), str),

            Ok(Signal::CtrlC) => continue,
            Ok(Signal::CtrlD) => break,

            Err(e) => {
                error!("Failed to read line\n{e}");
            }

            Ok(o) => {
                warn!("Unmanaged signal\n{:?}", o);
            }
        }
    }
}
