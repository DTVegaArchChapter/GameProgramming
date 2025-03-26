use piston_window;

mod draw;
mod game;
mod snake;

use piston_window::types::Color;
use piston_window::*;

use game::Game;

const BACK_COLOR: Color = [0.5, 0.5, 0.5, 1.0];

fn main() {
    let (width, height) = (30, 30);

    let mut window: PistonWindow =
        WindowSettings::new("Snake", [draw::to_coord(width), draw::to_coord(height)])
            .exit_on_esc(true)
            .build()
            .unwrap();

    let mut game = Game::new(width, height);

    while let Some(event) = window.next() {
        if let Some(Button::Keyboard(key)) = event.press_args() {
            game.key_pressed(key);
        }

        let mut glyphs = window.load_font("assets/Coolvetica.otf").unwrap();

        window.draw_2d(&event, |c, g, d| {
            clear(BACK_COLOR, g);
            game.draw(&c, g, &mut glyphs);
            glyphs.factory.encoder.flush(d);
        });

        event.update(|arg| {
            game.update(arg.dt);
        });
    }
}
