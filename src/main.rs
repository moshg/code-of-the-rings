use std::fmt;
use std::fmt::Write;
use std::io;
use std::io::BufRead;

/// 標準入力からマジックフレーズを読み込む.
pub fn input() -> Vec<u8> {
    let mut buf = Vec::with_capacity(501);
    let stdin = io::stdin();
    io::BufReader::new(stdin.lock()).read_until('\n' as u8, &mut buf).unwrap();
    buf[0..buf.len() - 1].iter().map(|&c| if c == ' ' as u8 { 0 } else { c - 'A' as u8 + 1 }).collect()
}

/// Bilboへの命令.
#[derive(Eq, PartialEq, Copy, Clone, Hash, Debug)]
pub enum Instr {
    Move(i8),
    Roll(i8),
    Trigger,
    Bra,
    Ket,
}

impl fmt::Display for Instr {
    fn fmt(&self, f: &mut fmt::Formatter) -> Result<(), fmt::Error> {
        match *self {
            Instr::Move(i) => {
                if i >= 0 {
                    for _ in 0..i {
                        f.write_char('>')?;
                    }
                } else {
                    for _ in 0..(-i) {
                        f.write_char('<')?;
                    }
                }
            }
            Instr::Roll(i) => {
                if i >= 0 {
                    for _ in 0..i {
                        f.write_char('+')?;
                    }
                } else {
                    for _ in 0..(-i) {
                        f.write_char('-')?;
                    }
                }
            }
            Instr::Trigger => f.write_char('.')?,
            Instr::Bra => f.write_char('[')?,
            Instr::Ket => f.write_char(']')?
        }
        Ok(())
    }
}

/// `from`から`to`までの距離を符号付きで求める.
///
/// `len`はトーラスの幅. `0 <= from, to < len`でなければならない.
#[inline]
fn dist_signed(from: u8, to: u8, len: u8) -> i8 {
    let (dist_right, dist_left) = if from < to {
        (to - from, len + from - to)
    } else {
        (len + to - from, from - to)
    };
    if dist_right <= dist_left { dist_right as i8 } else { -(dist_left as i8) }
}

/// `n`を`m`で割ったあまりを返す.
///
/// `-m < n < 2m - 1` でなければならない.
#[inline]
fn rem(n: i8, m: u8) -> u8 {
    (if n < 0 {
        n + m as i8
    } else if n >= m as i8 {
        n - m as i8
    } else {
        n
    }) as u8
}

/// マジックストーンの状態とBilboの位置.
#[derive(Eq, PartialEq, Copy, Clone, Default, Hash, Debug)]
pub struct State {
    stones: [u8; 30],
    pos: u8,
}

/// 命令列.
pub trait Instrs {
    /// 状態を更新する.
    fn update(&self, state: &mut State);

    /// 命令列をバッファに書き込む.
    fn write(&self, buf: &mut Vec<Instr>);

    /// 命令列の長さを返す.
    fn len(&self) -> u8;
}

// -27 < 内部データ < 27 でなければならない.
/// ルーンを変更する命令列.
#[derive(Eq, PartialEq, Copy, Clone, Hash, Debug)]
pub enum Roll {
    Roll(i8),
    ResetRoll(i8),
}

impl Default for Roll {
    #[inline]
    fn default() -> Roll {
        Roll::Roll(0)
    }
}

impl Instrs for Roll {
    #[inline]
    fn update(&self, state: &mut State) {
        let pos = state.pos as usize;
        let rune = match *self {
            Roll::Roll(n) => {
                let rune = state.stones[pos] as i8 + n;
                rem(rune, 27)
            }
            Roll::ResetRoll(n) => {
                rem(n, 27)
            }
        };
        state.stones[pos] = rune as u8;
    }

    #[inline]
    fn write(&self, buf: &mut Vec<Instr>) {
        match *self {
            Roll::Roll(n) => buf.push(Instr::Roll(n)),
            Roll::ResetRoll(n) => {
                buf.push(Instr::Bra);
                buf.push(Instr::Roll(1));
                buf.push(Instr::Ket);
                buf.push(Instr::Roll(n));
            }
        }
    }

    #[inline]
    fn len(&self) -> u8 {
        match *self {
            Roll::Roll(n) => n.abs() as u8,
            Roll::ResetRoll(n) => 3 + n.abs() as u8
        }
    }
}

/// ルーンを変更する効率的な命令列を求める.
///
/// # 探索候補
///
/// * 単純なロールの繰り返し
/// * ループでスペースにリセットした後, 単純なロールの繰り返し
#[inline]
pub fn calc_roll(stones: &[u8; 30], pos: u8, dest_rune: u8) -> Roll {
    let simple = dist_signed(stones[pos as usize], dest_rune, 27);
    let simple_cost = simple.abs();
    let reset = dist_signed(0, dest_rune, 27);
    let reset_cost = reset.abs() + 3;
    if simple_cost < reset_cost {
        Roll::Roll(simple)
    } else {
        Roll::ResetRoll(reset)
    }
}

#[derive(Eq, PartialEq, Copy, Clone, Hash, Debug)]
pub enum Move {
    Move(i8),
    RightSpaceMove(i8, u8),
    LeftSpaceMove(i8, u8),
}

impl Default for Move {
    #[inline]
    fn default() -> Self {
        Move::Move(0)
    }
}

impl Instrs for Move {
    #[inline]
    fn update(&self, state: &mut State) {
        match *self {
            Move::Move(n) => {
                let pos = state.pos as i8 + n;
                state.pos = rem(pos, 30);
            }
            Move::RightSpaceMove(_, pos) => {
                state.pos = pos;
            }
            Move::LeftSpaceMove(_, pos) => {
                state.pos = pos;
            }
        }
    }

    #[inline]
    fn write(&self, buf: &mut Vec<Instr>) {
        match *self {
            Move::Move(n) => buf.push(Instr::Move(n)),
            Move::RightSpaceMove(n, _) => {
                buf.push(Instr::Bra);
                buf.push(Instr::Move(1));
                buf.push(Instr::Ket);
                buf.push(Instr::Move(n));
            }
            Move::LeftSpaceMove(n, _) => {
                buf.push(Instr::Bra);
                buf.push(Instr::Move(-1));
                buf.push(Instr::Ket);
                buf.push(Instr::Move(n));
            }
        }
    }

    #[inline]
    fn len(&self) -> u8 {
        match *self {
            Move::Move(n) => n.abs() as u8,
            Move::RightSpaceMove(n, _) => 3 + n.abs() as u8,
            Move::LeftSpaceMove(n, _) => 3 + n.abs() as u8
        }
    }
}

/// 指定された位置の左右のスペースの位置を返す.
fn nearest_spaces(stones: &[u8; 30], pos: u8) -> Option<(u8, u8)> {
    if stones[pos as usize] == 0 {
        return Some((pos, pos));
    }

    let mut left = None;
    for offset in 1..30 {
        let pos = rem(pos as i8 - offset, 30);
        if stones[pos as usize] == 0 {
            left = Some(pos);
            break;
        }
    }

    // 左回りで0があるなら右回りでも0がある
    if let Some(left) = left {
        let mut right = 0;
        for offset in 1..30 {
            let pos = rem(pos as i8 + offset, 30);
            if stones[pos as usize] == 0 {
                right = pos;
                break;
            }
        }
        Some((left, right))
    } else {
        None
    }
}

/// 移動する効率的な命令列を求める.
#[inline]
pub fn calc_move(stones: &[u8; 30], from: u8, to: u8) -> Move {
    let simple = dist_signed(from, to, 30);
    let simple_len = simple.abs() as u8;

    if let Some((left, right)) = nearest_spaces(stones, from) {
        let left_move = dist_signed(left, to, 30);
        let left_len = left_move.abs() as u8 + 3;
        let right_move = dist_signed(right, to, 30);
        let right_len = right_move.abs() as u8 + 3;

        if right_len <= simple_len {
            if right_len <= left_len {
                Move::RightSpaceMove(right_move, to)
            } else {
                Move::LeftSpaceMove(left_move, to)
            }
        } else {
            if left_len <= simple_len {
                Move::LeftSpaceMove(left_move, to)
            } else {
                Move::Move(simple)
            }
        }
    } else {
        Move::Move(simple)
    }
}

/// 目的のルーン `rune` に対して, 現在の状態 `state` からそれを効率的にトリガーする命令列を`buf`にプッシュし,
/// `state`を更新する.
pub fn update(state: &mut State, buf: &mut Vec<Instr>, dest_rune: u8) {
    // 移動とロールの命令長の和を最小化する.
    let mut min_len = u8::max_value();
    let mut moving_min = Default::default();
    let mut rolling_min = Default::default();

    let pos = state.pos;
    for dest_pos in 0..30 {
        let moving = calc_move(&state.stones, pos, dest_pos);
        let rolling = calc_roll(&state.stones, dest_pos, dest_rune);

        let len = moving.len() + rolling.len();
        if len < min_len {
            min_len = len;
            moving_min = moving;
            rolling_min = rolling;
        }
    }

    // actionsを更新する.
    moving_min.write(buf);
    rolling_min.write(buf);
    buf.push(Instr::Trigger);

    // stateを更新する.
    moving_min.update(state);
    rolling_min.update(state);
}

fn main() {
    use std::io::Write;
    let input = input();
    let mut buf = Vec::new();
    let mut state = State::default();

    for &rune in &input {
        update(&mut state, &mut buf, rune);
    }

    let r = io::stdout();
    let mut r = io::BufWriter::new(r.lock());
    for instr in buf {
        write!(r, "{}", instr).unwrap();
    }
    write!(r, "\n").unwrap();
    r.flush().unwrap();
}

