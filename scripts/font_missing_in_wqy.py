from fontTools import ttLib

font = ttLib.ttFont.TTFont("../NotoSansCJKsc-Regular.otf")
font1 = ttLib.ttFont.TTFont("./stable/WenQuanYiZenHei-01.ttf")

def get_unicodes(f):
  a = []
  for ch, name in f["cmap"].getBestCmap().items():
    a.append(ch)

  return a

a = get_unicodes(font)
b = get_unicodes(font1)

c = []

for ch in a:
  if ch in b:
    continue
  c.append(ch)

for ch in c:
  print("U+{:4X} {}".format(ch, chr(ch)))
