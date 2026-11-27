// Dusk: logo, intro en achtergronden.
'use strict';

const FX = (() => {
  // ---------- pixel-logo ----------
  const GLYPHS = {
    D: ['####.', '#...#', '#...#', '#...#', '####.'],
    U: ['#...#', '#...#', '#...#', '#...#', '.###.'],
    S: ['.####', '#....', '.###.', '....#', '####.'],
    K: ['#...#', '#..#.', '###..', '#..#.', '#...#'],
  };
  const WORD = 'DUSK';

  function logoCells() {
    const cells = [];
    let x = 0;
    for (const ch of WORD) {
      const g = GLYPHS[ch];
      g.forEach((row, y) => [...row].forEach((c, gx) => { if (c === '#') cells.push([x + gx, y]); }));
      x += 6;
    }
    return { cells, w: x - 1, h: 5 };
  }

  // SVG-logo met schaduw, kleurverloop en een glans die eroverheen loopt.
  function logoSVG(size = 8, id = 'lg') {
    const { cells, w, h } = logoCells();
    const s = size, gap = Math.max(0.6, s * 0.07);
    const rects = (dx, dy, fill) => cells.map(([x, y]) =>
      `<rect x="${x * s + dx}" y="${y * s + dy}" width="${s - gap}" height="${s - gap}" rx="${s * 0.14}" fill="${fill}"/>`).join('');
    const W = (w + 1) * s + s * 0.4, H = h * s + s * 0.4;
    return `<svg class="pixlogo" viewBox="0 0 ${W} ${H}" width="${W}" height="${H}" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="${id}g" x1="0" x2="1" y1="0" y2="1">
          <stop offset="0" stop-color="#c084fc"/><stop offset=".55" stop-color="#a855f7"/><stop offset="1" stop-color="#f472b6"/>
        </linearGradient>
        <linearGradient id="${id}s" x1="0" x2="1">
          <stop offset="0" stop-color="#fff" stop-opacity="0"/><stop offset=".5" stop-color="#fff" stop-opacity=".75"/><stop offset="1" stop-color="#fff" stop-opacity="0"/>
        </linearGradient>
        <clipPath id="${id}c">${rects(0, 0, '#000')}</clipPath>
      </defs>
      <g opacity=".55">${rects(s * 0.4, s * 0.4, '#3b1670')}</g>
      <g>${rects(0, 0, `url(#${id}g)`)}</g>
      <g clip-path="url(#${id}c)"><g class="shine"><rect x="${-W * 0.45}" y="0" width="${W * 0.3}" height="${H}" fill="url(#${id}s)" transform="skewX(-20)"/></g></g>
    </svg>`;
  }

  // ---------- helpers ----------
  const hash = (x, y, s = 0) => {
    let h = (x * 374761393 + y * 668265263 + s * 2246822519) >>> 0;
    h = Math.imul(h ^ (h >>> 13), 1274126177) >>> 0;
    return ((h ^ (h >>> 16)) >>> 0) / 4294967295;
  };
  const sm = (a, b, x) => { const t = Math.min(1, Math.max(0, (x - a) / (b - a))); return t * t * (3 - 2 * t); };
  const RAMP = ' .:-=+*#%@';
  const rampCh = b => b <= 0 ? ' ' : RAMP[Math.min(RAMP.length - 1, Math.floor(b * (RAMP.length - 1)))];
  const mix = (a, b, t) => a.map((v, i) => Math.round(v + (b[i] - v) * Math.min(1, Math.max(0, t))));
  const rgba = (c, a = 1) => `rgba(${c[0]},${c[1]},${c[2]},${a})`;
  const PINK = [255, 95, 190], PURPLE = [168, 85, 247], VIOLET = [130, 80, 255], WHITE = [245, 238, 255];

  function fit(cv) {
    const dpr = Math.min(2, window.devicePixelRatio || 1);
    const w = cv.clientWidth, h = cv.clientHeight;
    if (cv.width !== Math.round(w * dpr) || cv.height !== Math.round(h * dpr)) {
      cv.width = Math.round(w * dpr); cv.height = Math.round(h * dpr);
    }
    const ctx = cv.getContext('2d');
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    return { ctx, w, h };
  }

  // ---------- intro: oog van tekst -> de ruimte in -> logo ----------
  function intro(cv, done) {
    const T_EYE = 3.4, T_LOGO = 1.9, TOTAL = T_EYE + T_LOGO;
    const start = performance.now();
    let raf, finished = false;
    const { cells, w: lw, h: lh } = logoCells();
    const parts = cells.map(([x, y], i) => ({ x, y, sx: hash(i, 1, 7), sy: hash(i, 2, 7), d: hash(i, 3, 7) }));

    function frame(now) {
      const t = (now - start) / 1000;
      const { ctx, w, h } = fit(cv);
      ctx.fillStyle = '#07060f';
      ctx.fillRect(0, 0, w, h);
      const cw = 9, ch = 16;
      ctx.font = `600 ${ch - 3}px "Cascadia Code", Consolas, monospace`;
      ctx.textBaseline = 'top';
      const cols = Math.ceil(w / cw), rows = Math.ceil(h / ch);

      if (t < T_EYE + 0.4) {
        const a = Math.min(cols / 2 * 0.78, rows * 1.75);
        let open = sm(0.05, 0.8, t);
        if (t > 1.35 && t < 1.6) open *= 1 - Math.sin((t - 1.35) / 0.25 * Math.PI) * 0.92;
        const look = t > 0.8 && t < 2.1 ? Math.sin((t - 0.8) / 1.3 * Math.PI * 2) * a * 0.16 : 0;
        const zt = sm(2.1, T_EYE, t);
        const zoom = Math.exp(zt * 3.6);
        const ri = a * 0.42, rp = ri * (0.36 + 0.26 * sm(1.8, 2.8, t));
        const fade = 1 - sm(T_EYE - 0.25, T_EYE + 0.4, t);
        for (let y = 0; y < rows; y++) for (let x = 0; x < cols; x++) {
          const px = (x - cols / 2) / zoom, py = (y - rows / 2) * 2.05 / zoom;
          const ix = px - look, d = Math.hypot(ix, py);
          const lid = a * 0.62 * open * (1 - (px / a) ** 2);
          let c = null, b = 0, chr;
          if (!(Math.abs(px) < a && Math.abs(py) < lid)) {
            const n = hash(x, y, 1) * 0.14 * (1 - zt);
            if (n > 0.06) { chr = rampCh(n * 3); c = [70, 58, 100]; b = 0.8; }
          } else if (d < rp) {
            const g = Math.exp(-((d / rp * 3.1) ** 2)) * (0.3 + 0.7 * zt);
            const st = hash(x, y, 77);
            if (st > 1 - 0.06 * zt) { const tw = 0.5 + 0.5 * Math.sin(t * 6 + st * 40); chr = tw > 0.7 ? '*' : '+'; c = WHITE; b = 0.4 + 0.6 * tw; }
            else if (g > 0.06) { chr = rampCh(g); c = mix([120, 20, 90], PINK, g); b = 1; }
          } else if (d < ri) {
            const th = Math.atan2(py, ix);
            const n = 0.5 + 0.25 * Math.sin(th * 23 + t * 0.6) + 0.25 * Math.sin(th * 41 + 1.3);
            const r = (d - rp) / (ri - rp);
            b = (0.45 + 0.5 * n * (1 - 0.4 * r)) * (r > 0.86 ? 0.55 : 1) * (r < 0.08 ? 0.6 : 1);
            chr = rampCh(b); c = mix(PINK, VIOLET, r); b = 0.55 + 0.45 * b;
          } else {
            const edge = 1 - Math.abs(py) / Math.max(lid, 0.001);
            if (edge < 0.1) { chr = '~'; c = PURPLE; b = 0.85; }
            else { const v = 0.3 + 0.35 * sm(0, 0.6, edge) - 0.15 * Math.abs(px) / a + 0.18 * (hash(x, y, 9) - 0.5); chr = rampCh(v); c = [225, 215, 240]; b = 0.45 + 0.4 * v; }
          }
          if (c && chr !== ' ') { ctx.fillStyle = rgba(c, b * fade); ctx.fillText(chr, x * cw, y * ch); }
        }
      }

      // logo vliegt in elkaar
      const lt = t - T_EYE + 0.15;
      if (lt > 0) {
        const s = Math.min(26, w / (lw + 8));
        const ox = (w - lw * s) / 2, oy = (h - lh * s) / 2 - s;
        const out = sm(T_LOGO - 0.45, T_LOGO, lt);
        for (const p of parts) {
          const k = sm(p.d * 0.35, 0.65 + p.d * 0.35, lt);
          const e = 1 - (1 - k) ** 3;
          const fx = (p.sx - 0.5) * w * 1.4 + w / 2, fy = (p.sy - 0.5) * h * 1.4 + h / 2;
          const x = fx + (ox + p.x * s - fx) * e, y = fy + (oy + p.y * s - fy) * e;
          const col = mix([192, 132, 252], [244, 114, 182], p.x / lw);
          ctx.globalAlpha = (0.2 + 0.8 * k) * (1 - out);
          ctx.shadowColor = rgba(col, 0.9); ctx.shadowBlur = 18 * k;
          ctx.fillStyle = rgba(col);
          const r = s * 0.86;
          ctx.beginPath(); ctx.roundRect(x, y, r, r, s * 0.18); ctx.fill();
        }
        ctx.shadowBlur = 0;
        ctx.globalAlpha = 1;
      }
      cv.parentElement.style.opacity = t > TOTAL - 0.35 ? String(1 - sm(TOTAL - 0.35, TOTAL, t)) : '1';
      if (t >= TOTAL) return finish();
      raf = requestAnimationFrame(frame);
    }
    function finish() {
      if (finished) return;
      finished = true;
      cancelAnimationFrame(raf);
      done();
    }
    raf = requestAnimationFrame(frame);
    return finish;
  }

  // ---------- achtergronden ----------
  const MODES = ['schemering', 'regen', 'golven', 'stil'];

  function background(cv) {
    let mode = 'schemering', calm = false, raf, last = 0;
    let stars = [], drops = [], meteors = [], W = 0, H = 0;

    function setup(w, h) {
      W = w; H = h;
      stars = Array.from({ length: Math.round(w * h / 5200) }, (_, i) => ({
        x: hash(i, 1) * w, y: Math.pow(hash(i, 2), 1.6) * h * 0.85, r: 0.4 + hash(i, 3) * 1.1,
        sp: 0.6 + hash(i, 4) * 2.2, ph: hash(i, 5) * 6.28, c: hash(i, 6),
      }));
      const cols = Math.ceil(w / 16);
      drops = Array.from({ length: cols }, (_, i) => ({ y: -hash(i, 9) * h * 1.5, v: 40 + hash(i, 10) * 110, len: 6 + Math.floor(hash(i, 11) * 18), s: i }));
    }

    function draw(now) {
      raf = requestAnimationFrame(draw);
      if (document.hidden) return;
      const fps = calm ? 12 : 30;
      if (now - last < 1000 / fps) return;
      const dt = Math.min(0.1, (now - last) / 1000 || 0.03);
      last = now;
      const t = now / 1000;
      const { ctx, w, h } = fit(cv);
      if (w !== W || h !== H) setup(w, h);
      ctx.clearRect(0, 0, w, h);
      if (mode === 'stil') return;

      if (mode === 'schemering') {
        for (const s of stars) {
          const tw = 0.55 + 0.45 * Math.sin(t * s.sp + s.ph);
          const fadeLow = 1 - sm(h * 0.45, h * 0.85, s.y);
          ctx.fillStyle = rgba(mix([200, 170, 255], WHITE, s.c), (0.25 + 0.6 * tw) * fadeLow);
          ctx.beginPath(); ctx.arc(s.x, s.y, s.r * (0.8 + 0.3 * tw), 0, 6.283); ctx.fill();
        }
        // vallende sterren
        if (!calm && Math.random() < dt * 0.12) meteors.push({ x: Math.random() * w * 0.8 + w * 0.2, y: Math.random() * h * 0.35, life: 0 });
        meteors = meteors.filter(m => (m.life += dt) < 1.1);
        for (const m of meteors) {
          const k = m.life / 1.1, x = m.x - k * 320, y = m.y + k * 140;
          const g = ctx.createLinearGradient(x, y, x + 120, y - 52);
          g.addColorStop(0, `rgba(255,235,255,${0.9 * (1 - k)})`); g.addColorStop(1, 'rgba(255,235,255,0)');
          ctx.strokeStyle = g; ctx.lineWidth = 1.6;
          ctx.beginPath(); ctx.moveTo(x, y); ctx.lineTo(x + 120, y - 52); ctx.stroke();
        }
        // wolken van tekst onderaan
        const cw = 11, ch = 16;
        ctx.font = `600 13px "Cascadia Code", Consolas, monospace`;
        ctx.textBaseline = 'top';
        const y0 = Math.floor(h * 0.6 / ch);
        for (let y = y0; y < Math.ceil(h / ch); y++) {
          const depth = (y - y0) / (h / ch - y0);
          for (let x = 0; x < Math.ceil(w / cw); x++) {
            const fx = x * 0.045 + t * (0.05 + depth * 0.08), fy = y * 0.16;
            let v = Math.sin(fx * 1.7 + Math.sin(fy * 1.3 + t * 0.07) * 1.4) + Math.sin(fx * 0.9 - fy * 0.8) * 0.8 + Math.sin(fx * 3.1 + fy * 2.2) * 0.35;
            v = (v + 2.15) / 4.3 * (0.35 + 0.75 * depth);
            if (v < 0.42) continue;
            const b = (v - 0.42) / 0.58;
            ctx.fillStyle = rgba(mix([110, 60, 190], [255, 120, 200], depth * 0.8 + b * 0.3), 0.08 + 0.32 * b);
            ctx.fillText(rampCh(Math.min(1, b * 0.9 + 0.1)), x * cw, y * ch);
          }
        }
      } else if (mode === 'regen') {
        ctx.font = `600 14px "Cascadia Code", Consolas, monospace`;
        ctx.textBaseline = 'top';
        const CH = '01<>/\\|=+*{}[]#$%&?!:;DUSK';
        for (const d of drops) {
          d.y += d.v * dt;
          if (d.y - d.len * 16 > h) { d.y = -Math.random() * h * 0.6; d.v = 40 + Math.random() * 110; }
          const x = d.s * 16;
          for (let j = 0; j < d.len; j++) {
            const y = Math.floor(d.y / 16) * 16 - j * 16;
            if (y < -16 || y > h) continue;
            const f = 1 - j / d.len;
            const c = j === 0 ? WHITE : mix([70, 30, 130], PURPLE, f);
            ctx.fillStyle = rgba(c, j === 0 ? 0.85 : 0.12 + 0.45 * f);
            ctx.fillText(CH[Math.floor(hash(d.s, y / 16, j === 0 ? Math.floor(t * 8) : 3) * CH.length)], x, y);
          }
        }
      } else if (mode === 'golven') {
        const cw = 11, ch = 16;
        ctx.font = `600 13px "Cascadia Code", Consolas, monospace`;
        ctx.textBaseline = 'top';
        for (let y = 0; y < Math.ceil(h / ch); y++) for (let x = 0; x < Math.ceil(w / cw); x++) {
          const fx = x, fy = y * 1.45;
          let v = Math.sin(fx * 0.06 + t * 0.8) + Math.sin(fy * 0.09 - t * 0.6) + Math.sin((fx + fy) * 0.04 + t * 0.45);
          v = (v + 3) / 6;
          if (v < 0.5) continue;
          const b = (v - 0.5) / 0.5;
          ctx.fillStyle = rgba(mix([90, 40, 170], PINK, b), 0.06 + 0.32 * b);
          ctx.fillText(rampCh(b * 0.85 + 0.1), x * cw, y * ch);
        }
      }
    }
    raf = requestAnimationFrame(draw);
    return {
      set(m) { mode = MODES.includes(m) ? m : 'schemering'; },
      calm(v) { calm = v; },
    };
  }

  return { logoSVG, intro, background, MODES };
})();
