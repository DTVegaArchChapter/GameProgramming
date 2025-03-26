use piston_window::types::Color;
use piston_window::*;

use rand::{thread_rng, Rng};

use crate::{
    draw::{self, draw_text},
    snake,
};

use draw::{draw_block, draw_rectangle};
use snake::{Direction, Snake};

const FOOD_COLOR: Color = [0.80, 0.00, 0.00, 1.0];
const BORDER_COLOR: Color = [0.00, 0.00, 0.00, 1.0];
const GAMEOVER_COLOR: Color = [0.80, 0.00, 0.00, 0.5];
const GAME_PAUSED_COLOR: Color = [0.00, 0.00, 0.00, 0.5];
const TEXT_WHITE: Color = [1.0, 1.0, 1.0, 1.0];

const MOVING_PERIOD: f64 = 0.1;
const RESTART_TIME: f64 = 1.0;

pub struct Game {
    snake: Snake,

    food_exists: bool,
    foot_x: i32,
    foot_y: i32,

    width: i32,
    height: i32,

    game_over: bool,
    waiting_time: f64,

    score: u32,

    game_paused: bool,
}

impl Game {
    pub fn new(width: i32, height: i32) -> Game {
        Game {
            snake: Snake::new(2, 2),
            food_exists: true,
            foot_x: 6,
            foot_y: 4,
            width,
            height,
            game_over: false,
            waiting_time: 0.0,
            score: 0,
            game_paused: false,
        }
    }

    pub fn key_pressed(&mut self, key: Key) {
        if self.game_over {
            return;
        }

        let dir = match key {
            Key::Up => Some(Direction::Up),
            Key::Down => Some(Direction::Down),
            Key::Left => Some(Direction::Left),
            Key::Right => Some(Direction::Right),
            Key::Space => {
                self.game_paused = !self.game_paused;
                return;
            }
            _ => None,
        };

        if dir.unwrap() == self.snake.head_direction().opposite() {
            return;
        }

        self.update_snake(dir);
    }

    pub fn draw(&self, con: &Context, g: &mut G2d, glyphs: &mut Glyphs) {
        self.snake.draw(con, g);

        if self.food_exists {
            draw_block(FOOD_COLOR, self.foot_x, self.foot_y, con, g);
        }

        draw_rectangle(BORDER_COLOR, 0, 0, self.width, 1, con, g);
        draw_rectangle(BORDER_COLOR, 0, self.height - 1, self.width, 1, con, g);
        draw_rectangle(BORDER_COLOR, 0, 0, 1, self.height, con, g);
        draw_rectangle(BORDER_COLOR, self.width - 1, 0, 1, self.height, con, g);

        let text = format!("Score: {}", self.score);

        draw_text(
            con,
            g,
            glyphs,
            TEXT_WHITE,
            [(self.width as u32), (self.height as u32)],
            &text,
        );

        if self.game_paused {
            draw_rectangle(GAME_PAUSED_COLOR, 0, 0, self.width, self.height, con, g);
            draw_text(
                con,
                g,
                glyphs,
                TEXT_WHITE,
                [
                    ((25 * self.width as u32 / 2) - 50),
                    (25 * self.height as u32 / 2),
                ],
                "Paused",
            );
        }

        if self.game_over {
            draw_rectangle(GAMEOVER_COLOR, 0, 0, self.width, self.height, con, g);
            draw_text(
                con,
                g,
                glyphs,
                TEXT_WHITE,
                [
                    ((25 * self.width as u32 / 2) - 50),
                    (25 * self.height as u32 / 2),
                ],
                "Game Over",
            );
            //drawing text for countdown
            let text = format!("Score: {}", self.score);
            draw_text(
                con,
                g,
                glyphs,
                TEXT_WHITE,
                [
                    ((25 * self.width as u32 / 2) - 50),
                    ((25 * self.height as u32 / 2) + 30),
                ],
                &text,
            );

            let text = format!("countdown: {:.2}", RESTART_TIME - self.waiting_time);
            draw_text(
                con,
                g,
                glyphs,
                TEXT_WHITE,
                [
                    ((25 * self.width as u32) - 200),
                    ((25 * self.height as u32) - 50),
                ],
                &text,
            );
        }
    }

    pub fn update(&mut self, delta_time: f64) {
        self.waiting_time += delta_time;

        if self.game_over {
            if self.waiting_time > RESTART_TIME {
                self.restart();
            }
            return;
        }

        if !self.food_exists {
            self.add_food();
        }

        if self.waiting_time > MOVING_PERIOD {
            self.update_snake(None);
        }

        if self.game_paused {
            return;
        }
    }

    pub fn chek_eating(&mut self) {
        let (head_x, head_y) = self.snake.head_position();

        if self.food_exists && self.foot_x == head_x && self.foot_y == head_y {
            self.food_exists = false;
            self.snake.restore_tail();
            self.score += 1;
        }
    }

    pub fn check_if_snake_alive(&mut self, dir: Option<Direction>) -> bool {
        let (next_x, next_y) = self.snake.next_head(dir);

        if self.snake.overlap_tail(next_x, next_y) {
            return false;
        }

        next_x > 0 && next_y > 0 && next_x < self.width - 1 && next_y < self.height - 1
    }

    pub fn add_food(&mut self) {
        let mut rng = thread_rng();

        let mut new_x = rng.gen_range(1..self.width - 1);
        let mut new_y = rng.gen_range(1..self.height - 1);

        while self.snake.overlap_tail(new_x, new_y) {
            new_x = rng.gen_range(1..self.width - 1);
            new_y = rng.gen_range(1..self.height - 1);
        }

        self.foot_x = new_x;
        self.foot_y = new_y;
        self.food_exists = true;
    }

    fn update_snake(&mut self, dir: Option<Direction>) {
        if self.game_paused {
            return;
        }

        if self.check_if_snake_alive(dir) {
            self.snake.move_forward(dir);
            self.chek_eating();
        } else {
            self.game_over = true;
        }

        self.waiting_time = 0.0;
    }

    fn restart(&mut self) {
        self.snake = Snake::new(2, 2);
        self.food_exists = true;
        self.foot_x = 6;
        self.foot_y = 4;
        self.game_over = false;
        self.waiting_time = 0.0;
    }
}
