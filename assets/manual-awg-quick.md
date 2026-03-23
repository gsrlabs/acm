```markdown
# How to properly use AmneziaWG 2.0 via console (awg-quick)

### 1. How to get the correct native config from the admin panel / client

1. Open AmneziaVPN (legacy or new) → select the AmneziaWG 2.0 server.
2. Click the gear icon ⚙️ next to the server.
3. Select **“Export configuration”** → **“AmneziaWG native format”** 
4. Save the file, for example as `awg-v2.conf`.

![Screenshot](screenshot_1.png)

**Important:**  
- Never use “Share” → the `vpn://...` string.  
- Use only the **native .conf**.

### 2. What must be edited in the AmneziaWG 2.0 config

The current version of `amneziawg-tools` **does not like** empty or zero values:
```ini
 I2= 
 I3= 
 I4= 
 I5=
```
**Simple rule:**

- Leave **only** the filled `I1=…`
- **Completely delete** the lines `I2=`, `I3=`, `I4=`, `I5=` (even if they contain `0` or just `I2 =`).

The corrected config should look something like this:

```ini
[Interface]
Address = 10.8.1.2/32
DNS = 1.1.1.1, 1.0.0.1
PrivateKey = ZRTNsTa60nhlg9ftfelvaeUehEvA2MP8nYHTkkmwI8vQ=
Jc = 4
Jmin = 10
Jmax = 50
S1 = 53
S2 = 67
S3 = 57
S4 = 15
H1 = 1169353282-1273681425
H2 = 1461293132-1822130420
H3 = 2109374387-2139701639
H4 = 2144623347-2145062022
I1 = <r 2><b 0x858000010001000000000669636c6f756403636f6d0000010001c00c000100010000105a00044d583737>

[Peer]
PublicKey = BYK7Sf3Sn4/HyGKjfkseFDLgXyntsmVx8j++yA=
PresharedKey = wFaNa3grsjkSA+efv/QrrcXfthjfrs9vIm8r0rK60=
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = 111.11.111.111:45000
PersistentKeepalive = 25
```

### 3. Commands to start without acm (cheat sheet)

```bash
# 1. Fix permissions
chmod 600 ~/awg-v2.conf

# 2. Bring up the connection
sudo awg-quick up ~/awg-v2.conf

# 3. Disconnect
sudo awg-quick down awg0
# or
sudo awg-quick down ~/awg-v2.conf

# 4. Check status
awg show

# 5. If you placed it in the standard location
sudo cp ~/awg-v2.conf /etc/wireguard/awg0.conf
chmod 600 /etc/wireguard/awg0.conf
sudo awg-quick up awg0

```

### 4. Autostart on boot

```bash
sudo systemctl enable --now awg-quick@wg0.service
```

### 5. Useful tips

- After any config change — do `down + up`.
- `I1` is the main obfuscation parameter. If it is present and correct — the protocol works in 2.0 mode.
- `I2`–`I5` can be omitted entirely (they are optional). If you want to fill them in — you need real CPS packets in the format `<b 0x...>` or `<r X><b 0x...>`.
```