# Code of the Rings

My solution to the Code of the Rings challenge.

https://www.codingame.com/multiplayer/optimization/code-of-the-rings

## 考察

トリガーの方針
1. まず両隣のどちらかにいく

2. ループ or 愚直にトリガー

2.a. ループの場合
ループカウンタを最大値 or 残り回数にセットする
ループ内ではトリガーしながら進んだあと戻る、または進んだあとトリガーしながら戻る

2文字を2回トリガーする例: `++[>.>.<<-]`

スペル数の考察:
* 文字数: l
* トリガー回数と現在のカウンタの距離: offset
offset + 2 + 2l + l + 1 = 3 + 3l + offset

2.b. ループしない場合
「トリガーしながら進んだあと戻る、または進んだあとトリガーしながら戻る」を繰り返す

2文字を2回トリガーする例: `.>.<.>.`

スペル数の考察:
* 文字数: l
* トリガー回数: n
l + (l - 1) + (n - 1) * (l + (l - 1) * 2) = (3l - 2)n - l + 1
