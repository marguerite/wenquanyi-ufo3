WenQuanYi Zen Hei in Unified Font Object 3 format, and can be built with googlefonts' fontmake.

Based on the latest snapshot version in 20150916

Changes:

20240209:
  * 重新基于 0.9.45 稳定版本生成 Medium，因为通过 fontforge 的 sfddiff 比较了 0.9.47 snapshot，
    除了 U+3000 区域的少数符号外（snapshot 里有的还都处于不可用状态），真正的汉字只有 U+9FCC 鿌是新增的。
  * 与 Noto Sans CJK SC 相比生成了待补字清单。后续准备通过 zi2zi-pytorch 补齐
