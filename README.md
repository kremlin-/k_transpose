# kremlin tranpose
*top secret kremlin color transposition device*

transposes colors in web pages to match your ANSI terminal color theme.

k_transpose starts a local http server (localhost/kt/) through which you browse the outside web in the form e.g. http://localhost/kt/google.com. both relative & explicit URLs/URIs are prepended to route through this local http server such that one can browse the web "normally" after specifying the transposition server once

terminal colors are read from ~/.Xresources in the form "URxvt*color[0-16]:#RRGGBB or *color[0-16]:#RRGGBB. *background: & *foreground: colors are used too, to determine light/dark theme context

##colorspace transposition
TODO

### compilation/installation
```
make
sudo make install
```

### running
```
# you likely need root to bind httpds to a local port
sudo k_transpose
```
