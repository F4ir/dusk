// Dusk: de app zelf.
'use strict';

const $ = s => document.querySelector(s);
const esc = s => String(s ?? '').replace(/[&<>"']/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));

// ---------- iconen ----------
const P = {
  home: '<path d="M3 10.5 12 3l9 7.5V20a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1z"/>',
  apps: '<rect x="3" y="3" width="7.5" height="7.5" rx="2.2"/><rect x="13.5" y="3" width="7.5" height="7.5" rx="2.2"/><rect x="3" y="13.5" width="7.5" height="7.5" rx="2.2"/><rect x="13.5" y="13.5" width="7.5" height="7.5" rx="2.2"/>',
  tweaks: '<path d="M4 6h9M17 6h3M4 12h3M11 12h9M4 18h11M19 18h1"/><circle cx="15" cy="6" r="2"/><circle cx="9" cy="12" r="2"/><circle cx="17" cy="18" r="2"/>',
  clean: '<path d="M11 3.5l1.9 5.1 5.1 1.9-5.1 1.9L11 17.5l-1.9-5.1L4 10.5l5.1-1.9z"/><path d="M18.5 15l.8 2.2 2.2.8-2.2.8-.8 2.2-.8-2.2-2.2-.8 2.2-.8z"/>',
  settings: '<circle cx="12" cy="12" r="3.2"/><path d="M19.4 15a1.7 1.7 0 0 0 .3 1.8l.1.1a2 2 0 1 1-2.8 2.8l-.1-.1a1.7 1.7 0 0 0-1.8-.3 1.7 1.7 0 0 0-1 1.5V21a2 2 0 1 1-4 0v-.1a1.7 1.7 0 0 0-1.1-1.5 1.7 1.7 0 0 0-1.8.3l-.1.1a2 2 0 1 1-2.8-2.8l.1-.1a1.7 1.7 0 0 0 .3-1.8 1.7 1.7 0 0 0-1.5-1H3a2 2 0 1 1 0-4h.1a1.7 1.7 0 0 0 1.5-1.1 1.7 1.7 0 0 0-.3-1.8l-.1-.1a2 2 0 1 1 2.8-2.8l.1.1a1.7 1.7 0 0 0 1.8.3H9a1.7 1.7 0 0 0 1-1.5V3a2 2 0 1 1 4 0v.1a1.7 1.7 0 0 0 1 1.5 1.7 1.7 0 0 0 1.8-.3l.1-.1a2 2 0 1 1 2.8 2.8l-.1.1a1.7 1.7 0 0 0-.3 1.8V9a1.7 1.7 0 0 0 1.5 1H21a2 2 0 1 1 0 4h-.1a1.7 1.7 0 0 0-1.5 1z"/>',
  refresh: '<path d="M20 11.5A8 8 0 1 1 17.7 6M20 4v5h-5"/>',
  search: '<circle cx="11" cy="11" r="7"/><path d="m20 20-3.5-3.5"/>',
  check: '<path d="m5 12.5 4.5 4.5L19 7.5"/>',
  down: '<path d="M12 4v11m0 0-4.5-4.5M12 15l4.5-4.5M5 20h14"/>',
  up: '<path d="M12 19V6m0 0-5 5m5-5 5 5"/>',
  shield: '<path d="M12 3 4.5 6v6c0 4.4 3.2 7.9 7.5 9 4.3-1.1 7.5-4.6 7.5-9V6z"/><path d="m9 12 2 2 4-4"/>',
  lock: '<rect x="5" y="11" width="14" height="10" rx="2.5"/><path d="M8 11V8a4 4 0 0 1 8 0v3"/>',
  trash: '<path d="M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7V4h6v3"/>',
  pkg: '<path d="M12 3 3.5 7.5v9L12 21l8.5-4.5v-9z"/><path d="M3.5 7.5 12 12l8.5-4.5M12 12v9"/>',
  folder: '<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>',
  palette: '<path d="M12 3a9 9 0 1 0 0 18c1.1 0 1.6-.9 1.2-1.8-.5-1.1.2-2.2 1.4-2.2H17a4 4 0 0 0 4-4c0-5.5-4-10-9-10z"/><circle cx="7.5" cy="11" r="1.2"/><circle cx="10.5" cy="7" r="1.2"/><circle cx="15" cy="7.5" r="1.2"/>',
  game: '<rect x="2" y="6" width="20" height="12" rx="6"/><path d="M7 10v4M5 12h4M15 11h.01M18 13h.01"/>',
  zap: '<path d="M13 2 4 14h7l-1 8 9-12h-7z"/>',
  eyeoff: '<path d="M3 3l18 18M10.6 6.1c.5-.1.9-.1 1.4-.1 5 0 8.5 4 9.5 6-.4.9-1.3 2.2-2.6 3.5M6.6 7.6C4.6 8.9 3.2 10.7 2.5 12c1 2 4.5 6 9.5 6 1.6 0 3-.4 4.3-1M9.9 9.9a3 3 0 0 0 4.2 4.2"/>',
  x: '<path d="M6 6l12 12M18 6 6 18"/>',
  alert: '<path d="M12 9v4M12 17h.01"/><path d="M10.3 3.9 2.4 18a2 2 0 0 0 1.7 3h15.8a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0z"/>',
  spark: '<path d="M12 3l2 6 6 2-6 2-2 6-2-6-6-2 6-2z"/>',
  cpu: '<rect x="6" y="6" width="12" height="12" rx="2"/><path d="M9 2v4M15 2v4M9 18v4M15 18v4M2 9h4M2 15h4M18 9h4M18 15h4"/>',
  play: '<path d="M7 4.5v15l13-7.5z"/>',
};
const ic = (n, cls = '') => `<svg class="${cls}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">${P[n]}</svg>`;
const MOON = '<svg viewBox="0 0 24 24"><defs><linearGradient id="mg2" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#c084fc"/><stop offset="1" stop-color="#f472b6"/></linearGradient></defs><path fill="url(#mg2)" d="M20.5 14.6A8.6 8.6 0 0 1 9.4 3.5a8.6 8.6 0 1 0 11.1 11.1Z"/><circle cx="17" cy="5" r=".9" fill="#fff"/><circle cx="20.5" cy="8.5" r=".6" fill="#fff" opacity=".7"/></svg>';

P.laptop = '<rect x="4" y="5" width="16" height="11" rx="2"/><path d="M2 19h20"/>';
P.desktop = '<rect x="3" y="4" width="18" height="12" rx="2"/><path d="M9 20h6M12 16v4"/>';
P.battery = '<rect x="2" y="7" width="18" height="10" rx="2"/><path d="M22 11v2M6 10v4M10 10v4"/>';
P.rocket = '<path d="M5 15c-1.5 1.5-2 5-2 5s3.5-.5 5-2M9 12a22 22 0 0 1 3-6c2.5-3 6-3.5 9-3-.5 3-1 6.5-4 9a22 22 0 0 1-6 3z"/><circle cx="15" cy="9" r="1.6"/><path d="M9 12H5l2-4h4M12 15v4l4-2v-4"/>';
P.wrench = '<path d="M14.7 6.3a4 4 0 0 0 5 5L21 13l-8 8-3-3 1.3-1.3a4 4 0 0 0-5-5L3 10l8-8 3 3z"/>';
P.chip = '<rect x="5" y="5" width="14" height="14" rx="2.5"/><rect x="9" y="9" width="6" height="6" rx="1"/><path d="M9 2v3M15 2v3M9 19v3M15 19v3M2 9h3M2 15h3M19 9h3M19 15h3"/>';
P.globe = '<circle cx="12" cy="12" r="9"/><path d="M3 12h18M12 3c2.5 2.6 3.8 5.6 3.8 9s-1.3 6.4-3.8 9c-2.5-2.6-3.8-5.6-3.8-9S9.5 5.6 12 3z"/>';
P.ext = '<path d="M14 4h6v6M20 4l-9 9M18 14v5a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1h5"/>';
P.disk = '<rect x="3" y="6" width="18" height="12" rx="3"/><path d="M7 14h.01M11 14h6"/>';
P.board = '<rect x="3" y="3" width="18" height="18" rx="3"/><path d="M8 8h3v3H8zM14 8h2M14 11h2M8 15h8"/>';
P.gpu = '<rect x="2" y="7" width="18" height="10" rx="2"/><circle cx="9" cy="12" r="2.5"/><circle cx="15" cy="12" r="2.5"/><path d="M20 10h2M6 17v3"/>';
P.ram = '<rect x="3" y="7" width="18" height="9" rx="1.5"/><path d="M7 10v3M11 10v3M15 10v3M5 16v2M19 16v2"/>';

// ---------- koppeling met Go ----------
const API = () => window.go.main.App;
const on = (name, fn) => window.runtime.EventsOn(name, fn);

// ---------- staat ----------
const S = {
  page: 'home', sys: null,
  apps: [], appsLoaded: false, appsErr: null, scanning: false,
  filter: 'all', view: 'all', q: '', progs: null, progSel: new Set(), progProg: {}, progSort: 'name', sel: new Set(), prog: {},
  tweaks: [], tweaksLoaded: false, tweakWait: new Set(), explorerPending: false,
  bloat: [], bloatLoaded: false, bloatErr: null, bloatSel: new Set(), bloatProg: {},
  hw: null, hwErr: null, drv: null, drvScanning: false, drvErr: null, drvSel: new Set(), drvProg: {},
  busy: false, feed: [], failed: new Set(), twOpen: new Set(), showIgnored: false,
  drvChecking: false, drvStat: null, rpProg: {}, pv: null, pvOpen: new Set(), ctx: null, upd: null, updPct: null, gpuProg: {}, oemToolProg: {}, oemSel: new Set(), oemProg: {}, st: null, junk: null, junkSel: new Set(['temp', 'shader', 'browser', 'discord', 'crash', 'wintemp', 'wucache']), junkBusy: false,
  set: Object.assign({ bg: 'schemering', intro: true, calm: false, lang: 'en', ignored: [] }, JSON.parse(localStorage.getItem('dusk') || '{}')),
};
const save = () => localStorage.setItem('dusk', JSON.stringify(S.set));
LANG = S.set.lang;

const PAGES = [
  { id: 'home', icon: 'home' },
  { id: 'apps', icon: 'apps', group: 'software' }, { id: 'startup', icon: 'rocket', group: 'software' }, { id: 'clean', icon: 'clean', group: 'software' },
  { id: 'tweaks', icon: 'tweaks', group: 'system' }, { id: 'privacy', icon: 'eyeoff', group: 'system' },
  { id: 'hardware', icon: 'chip', group: 'system' }, { id: 'repair', icon: 'wrench', group: 'system' },
  { id: 'settings', icon: 'settings', bottom: true },
];
const CATS = ['all', 'popular', 'installed', 'updates', 'oem', 'browser', 'chat', 'ai', 'games', 'gear', 'media', 'tools', 'system', 'network', 'dev', 'runtime'];
const GROUP_ICON = { looks: 'palette', taskbar: 'apps', explorer: 'folder', gaming: 'game', speed: 'zap', privacy: 'eyeoff', system: 'settings' };
const GROUP_ORDER = ['looks', 'taskbar', 'explorer', 'gaming', 'speed', 'system', 'privacy'];

// ---------- score ----------
const RANKS = [{ min: 90 }, { min: 70 }, { min: 45 }, { min: 0 }];
const rankName = min => t('rank_' + min);

function scoreParts() {
  const parts = [];
  if (S.appsLoaded) {
    const inst = S.apps.filter(a => a.Installed), upd = inst.filter(a => a.Available).length;
    parts.push({
      icon: 'up', name: t('sp_apps'), max: 40, got: 40 * (inst.length ? (inst.length - upd) / inst.length : 1),
      text: t('sp_apps_text', inst.length - upd, inst.length), todo: upd ? { label: t('update_all'), act: 'updAll' } : null,
    });
  }
  // drivers in plaats van tweaks: tweaks zijn een eigen keuze, drivers niet
  if (S.hw && S.drvStat) {
    const ign = S.set.ignored || [];
    const bad = (S.hw.Devices || []).filter(d => !ign.includes(d.Name.trim()) && ['update', 'missing', 'problem'].includes(devStatus(d).k)).length + (S.drvStat.OEM || []).length;
    parts.push({
      icon: 'chip', name: t('sp_drivers'), max: 35, got: 35 * Math.max(0, 1 - bad * 0.15),
      text: t('sp_drivers_text', bad), todo: bad ? { label: t('sp_drivers_btn'), go: 'hardware' } : null,
    });
  }
  if (S.bloatLoaded) {
    const n = S.bloat.length;
    parts.push({
      icon: 'clean', name: t('sp_clean'), max: 25, got: 25 * Math.max(0, 1 - n * 0.12),
      text: t('sp_clean_text', n), todo: n ? { label: t('clean_up'), go: 'clean' } : null,
    });
  }
  return parts;
}

function score() {
  const parts = scoreParts();
  const max = parts.reduce((a, p) => a + p.max, 0);
  if (!max) return { pct: 0, rank: t('scanning'), ready: false };
  const pct = Math.round(parts.reduce((a, p) => a + p.got, 0) / max * 100);
  return { pct, rank: rankName(RANKS.find(r => pct >= r.min).min), ready: true };
}

function showScore() {
  const sc = score(), parts = scoreParts();
  const m = $('#modal');
  m.innerHTML = `<div class="dialog score-dlg">
    <button class="btn ghost sm icon dlg-x" data-close title="${t('close')}">${ic('x')}</button>
    <div class="sd-top"><svg class="ring" viewBox="0 0 96 96" width="96" height="96">${ring(96, 8, sc.pct, 'sdr')}</svg>
      <div><div class="sd-pct">${sc.pct}%</div><div class="sd-rank">${esc(sc.rank)}</div></div></div>
    <p class="sd-intro">${t('score_intro')}</p>
    <div class="sd-parts">${parts.map(p => `
      <div class="sd-part">
        <div class="sd-ic">${ic(p.icon)}</div>
        <div class="sd-body">
          <div class="sd-row"><b>${p.name}</b><span>${t('points', Math.round(p.got), p.max)}</span></div>
          <div class="bar"><i style="width:${p.got / p.max * 100}%"></i></div>
          <div class="sd-text">${esc(p.text)}</div>
        </div>
        ${p.todo ? `<button class="btn sm" ${p.todo.act ? `data-act="${p.todo.act}"` : `data-goto="${p.todo.go}"`} data-close>${esc(p.todo.label)}</button>` : `<span class="pill ok">${ic('check').replace('<svg', '<svg width="13" height="13"')}${t('top')}</span>`}
      </div>`).join('')}</div>
    <div class="sd-ranks">${[...RANKS].reverse().map((r, i, a) => `<div class="${rankName(r.min) === sc.rank ? 'on' : ''}"><b>${rankName(r.min)}</b><span>${r.min}${i < a.length - 1 ? '-' + (a[i + 1].min - 1) : '+'}%</span></div>`).join('')}</div>
    <p class="sd-note">${t('score_note')}</p>
  </div>`;
  m.classList.add('on');
  const close = () => { m.classList.remove('on'); m.onclick = null; document.removeEventListener('keydown', key); };
  const key = e => { if (e.key === 'Escape') close(); };
  document.addEventListener('keydown', key);
  m.onclick = e => { if (e.target === m || e.target.closest('[data-close]')) setTimeout(close, 0); };
}

function ring(size, sw, pct, id) {
  const r = (size - sw) / 2, c = 2 * Math.PI * r;
  return `<defs><linearGradient id="${id}" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#a855f7"/><stop offset="1" stop-color="#f472b6"/></linearGradient></defs>
    <circle cx="${size / 2}" cy="${size / 2}" r="${r}" fill="none" stroke="rgba(255,255,255,.08)" stroke-width="${sw}"/>
    <circle cx="${size / 2}" cy="${size / 2}" r="${r}" fill="none" stroke="url(#${id})" stroke-width="${sw}" stroke-linecap="round"
      stroke-dasharray="${c}" stroke-dashoffset="${c * (1 - pct / 100)}" transform="rotate(-90 ${size / 2} ${size / 2})" style="filter:drop-shadow(0 0 6px rgba(236,72,153,.6))"/>`;
}

// ---------- algemeen ----------
function toast(text, kind = '') {
  const el = document.createElement('div');
  el.className = 'toast ' + kind;
  el.innerHTML = `<div class="ic">${ic(kind === 'ok' ? 'check' : kind === 'bad' ? 'alert' : 'spark')}</div><div>${esc(text)}</div>`;
  $('#toasts').append(el);
  setTimeout(() => { el.classList.add('out'); setTimeout(() => el.remove(), 300); }, 4200);
}

function confirmBox({ title, text, ok, danger = false }) {
  return new Promise(res => {
    const m = $('#modal');
    m.innerHTML = `<div class="dialog"><h3>${esc(title)}</h3><p>${text}</p><div class="acts">
      <button class="btn ghost" data-v="0">${t('cancel')}</button><button class="btn ${danger ? 'danger' : 'primary'}" data-v="1">${esc(ok)}</button></div></div>`;
    m.classList.add('on');
    const close = v => { m.classList.remove('on'); m.onclick = null; document.removeEventListener('keydown', key); res(v); };
    const key = e => { if (e.key === 'Escape') close(false); if (e.key === 'Enter') close(true); };
    document.addEventListener('keydown', key);
    m.onclick = e => { const b = e.target.closest('[data-v]'); if (b) close(b.dataset.v === '1'); else if (e.target === m) close(false); };
  });
}

const locale = () => LANG === 'nl' ? 'nl-NL' : 'en-GB';
const now = () => new Date().toLocaleTimeString(locale(), { hour: '2-digit', minute: '2-digit' });

function log(line) {
  const el = document.createElement('div');
  const s = String(line);
  if (s.startsWith('»')) {
    el.className = 'h';
    el.innerHTML = `<span class="ts">${now()}</span>${esc(s.slice(2))}`;
    S.feed.unshift({ t: now(), text: s.slice(2) });
    S.feed = S.feed.slice(0, 6);
    if (S.page === 'home') renderHome();
  } else {
    el.textContent = s;
    if (s.includes('✓')) el.className = 'ok';
    if (s.includes('✗')) el.className = 'bad';
  }
  const box = $('#log');
  const stick = box.scrollTop + box.clientHeight >= box.scrollHeight - 30;
  box.append(el);
  while (box.children.length > 600) box.firstChild.remove();
  if (stick) box.scrollTop = box.scrollHeight;
}

function setBusy(b) {
  S.busy = b;
  $('#btnLog').classList.toggle('busy', b);
}

// vaste teksten in de HTML (data-t)
function applyStatic() {
  document.documentElement.lang = LANG;
  document.querySelectorAll('[data-t]').forEach(el => { el.textContent = t(el.dataset.t); });
}

// ---------- navigatie ----------
function renderNav() {
  const upd = S.apps.filter(a => a.Available).length;
  let probs = S.hw ? visibleProblems().length : 0;
  if (S.hw && S.drvStat) probs = (S.hw.Devices || []).filter(d => !(S.set.ignored || []).includes(d.Name.trim()) && ['update', 'missing', 'problem'].includes(devStatus(d).k)).length + (S.drvStat.OEM || []).length;
  const btn = p => {
    let badge = '';
    if (p.id === 'apps' && upd) badge = `<span class="badge">${upd}</span>`;
    if (p.id === 'clean' && S.bloat.length) badge = `<span class="badge pink">${S.bloat.length}</span>`;
    if (p.id === 'hardware' && probs) badge = `<span class="badge pink">${probs}</span>`;
    return `<button class="nav ${S.page === p.id ? 'on' : ''}" data-page="${p.id}">${ic(p.icon)}<span>${t('nav_' + p.id)}</span>${badge}</button>`;
  };
  // groepjes met een kopje; Instellingen los onderaan
  let html = '', grp = '';
  for (const p of PAGES.filter(x => !x.bottom)) {
    if (p.group && p.group !== grp) { grp = p.group; html += `<div class="nav-h">${t('nav_grp_' + grp)}</div>`; }
    html += btn(p);
  }
  $('#nav').innerHTML = html;
  $('#navBottom').innerHTML = PAGES.filter(x => x.bottom).map(btn).join('');
  const sc = score();
  $('#miniRing').innerHTML = ring(44, 5, sc.pct, 'mr');
  $('#miniRank').textContent = sc.rank;
  $('#miniPct').textContent = sc.ready ? `${sc.pct}% · ${t('pc_score')}` : t('pc_score');
  $('.mini-score').title = t('what_means');
}

function goPage(page) {
  S.page = page;
  document.querySelectorAll('.page').forEach(p => p.classList.toggle('on', p.id === 'p-' + page));
  renderNav();
  render();
  if (page === 'hardware' && S.hw && !S.drvStat && !S.drvChecking) checkDrivers();
  if (page === 'hardware' && !S.hw && !S.hwErr) loadHardware().then(() => { if (!S.drvStat) checkDrivers(); });
  if (page === 'startup' && !S.st) loadStartup();
  if (page === 'clean' && !S.junk) scanJunk();
  if (page === 'clean' && !S.ctx) loadCtx();
  if (page === 'privacy') loadPrivacy();
}

function render() {
  ({ home: renderHome, apps: renderApps, tweaks: renderTweaks, startup: renderStartup, hardware: renderHardware, clean: renderClean, settings: renderSettings, privacy: renderPrivacy, repair: renderRepair })[S.page]();
  renderSel();
}

// Maakt tekst kleiner tot hij past; past hij dan nog niet, dan op meer regels.
function fitText(el, max, min) {
  el.classList.remove('wrap');
  let size = max;
  el.style.fontSize = size + 'px';
  while (el.scrollWidth > el.clientWidth && size > min) {
    size -= 0.5;
    el.style.fontSize = size + 'px';
  }
  if (el.scrollWidth > el.clientWidth) el.classList.add('wrap');
}

// ---------- overzicht ----------
function greet() {
  const h = new Date().getHours();
  return t(h < 6 ? 'greet_night' : h < 12 ? 'greet_morning' : h < 18 ? 'greet_afternoon' : 'greet_evening');
}

function renderHome() {
  const s = S.sys || {};
  const sc = score();
  const upd = S.apps.filter(a => a.Available);
  const twOn = S.tweaks.filter(x => x.On).length;
  const name = (s.User || '').toUpperCase();
  $('#p-home').innerHTML = `${updBanner()}
    <div class="card hero">
      <div>
        <div class="kick">${greet()},</div>
        <h2 title="${esc(name)}" style="font-size:${name.length <= 12 ? 38 : name.length <= 20 ? 30 : 24}px"><em>${esc(name || t('welcome'))}</em></h2>
        <div class="meta">
          <span class="pill pink" title="${esc(s.PC)}">${esc(s.PC || '…')}</span>
          <span class="pill" title="${esc(s.Windows)}">${esc(s.Windows || '')}</span>
          ${s.Admin ? `<span class="pill ok"><i></i>${t('admin')}</span>` : `<span class="pill adm">${t('normal_user')}</span>`}
        </div>
      </div>
      <div class="ring-big" data-score title="${t('what_means')}"><svg class="ring" viewBox="0 0 150 150" width="150" height="150">${ring(150, 11, sc.pct, 'hr')}</svg>
        <div class="val"><b>${sc.ready ? sc.pct + '%' : '…'}</b><span>${esc(sc.rank)}</span></div><div class="why">${ic('spark')}${t('what_means')}</div></div>
    </div>
    <div class="stats">
      <div class="card stat">
        <div class="ic" style="background:rgba(251,191,36,.13);color:var(--warn)">${ic('up')}</div>
        <div class="num">${S.appsLoaded ? upd.length : '<span class="skel" style="width:40px;height:28px"></span>'}</div>
        <div class="lbl">${t('updates_ready', upd.length)}</div>
        ${upd.length ? `<button class="btn primary sm" data-act="updAll">${ic('down')}${t('update_all')}</button>` : `<button class="btn sm" data-goto="apps">${t('to_apps')}</button>`}
      </div>
      <div class="card stat">
        <div class="ic" style="background:rgba(168,85,247,.16);color:#d8b4fe">${ic('tweaks')}</div>
        <div class="num">${S.tweaksLoaded ? `${twOn}<span style="font-size:18px;color:var(--faint)"> / ${S.tweaks.length}</span>` : '…'}</div>
        <div class="lbl">${t('tweaks_on')}</div>
        <button class="btn sm" data-goto="tweaks">${t('see_tweaks')}</button>
      </div>
      <div class="card stat">
        <div class="ic" style="background:rgba(236,72,153,.15);color:#ff9ccc">${ic('clean')}</div>
        <div class="num">${S.bloatLoaded ? S.bloat.length : '…'}</div>
        <div class="lbl">${t('junk_found')}</div>
        <button class="btn sm" data-goto="clean">${S.bloat.length ? t('clean_up') : t('view')}</button>
      </div>
    </div>
    <div class="section-title">${t('system')}</div>
    <div class="card specs">
      <div class="spec"><div class="k">${t('processor')}</div><div class="v" title="${esc(s.CPU)}">${esc(shortCPU(s.CPU) || '…')}</div></div>
      <div class="spec"><div class="k">${t('memory')}</div><div class="v">${s.RAM ? s.RAM + ' GB' : '…'}</div></div>
      <div class="spec"><div class="k">${t('windows')}</div><div class="v">${esc(s.Windows || '…')}</div></div>
      <div class="spec"><div class="k">${t('winget')}</div><div class="v" style="color:${s.Winget === false ? 'var(--bad)' : 'var(--ok)'}">${s.Winget === false ? t('not_found') : t('ready')}</div></div>
    </div>
    <div class="section-title">${t('recent_activity')}</div>
    <div class="card feed">${S.feed.length ? S.feed.map(f => `<div class="row"><span class="d"></span>${esc(f.text)}<span class="t">${f.t}</span></div>`).join('')
      : `<div class="empty">${t('nothing_yet')}</div>`}</div>`;
}

const shortCPU = c => (c || '').replace(/\s+/g, ' ').replace(/\s*\d+-Core Processor/i, '').replace(/\(R\)|\(TM\)/g, '').trim();

// ---------- apps ----------
function appIcon(a) {
  const inner = a.Icon ? `<i style="--m:url('icons/${a.Icon}.svg')"></i>` : esc(a.Name.replace(/[^A-Za-z0-9]/g, '').slice(0, 2));
  return `<div class="ico" style="--c:${a.Color}">${inner}</div>`;
}

function stepText(p) { return p ? t('step_' + p.step) : ''; }

function appCard(a) {
  const p = S.prog[a.ID];
  const busy = p && (p.phase === 'wacht' || p.phase === 'bezig');
  const sel = S.sel.has(a.ID);
  let status, btn = '', dot = '';
  if (!S.appsLoaded) status = '<span class="skel" style="width:90px;height:16px"></span>';
  else if (!a.Installed) {
    status = `<span class="st-txt">${t('not_installed')}</span>`;
    btn = `<button class="btn sm" data-act="one" data-id="${esc(a.ID)}">${ic('down')}${t('install')}</button>`;
  } else if (a.Available) {
    dot = 'upd';
    status = `<span class="st-txt upd" title="${esc(a.Version)} → ${esc(a.Available)}">${t('upd_pill')} · ${esc(a.Available)}</span>`;
    btn = `<button class="btn primary sm" data-act="one" data-id="${esc(a.ID)}">${ic('up')}${t('update')}</button>`;
  } else {
    dot = 'ok';
    status = `<span class="st-txt ok" title="${esc(a.Version)}">${ic('check')}${t('installed_v')}${a.Version && !/unknown/i.test(a.Version) ? ' · ' + esc(a.Version) : ''}</span>`;
  }
  const pct = p && p.pct >= 0 ? p.pct : null;
  return `<div class="card app ${sel ? 'sel' : ''} ${busy ? 'busy' : ''}" data-id="${esc(a.ID)}">
    <div class="app-row"><div class="ico-wrap">${appIcon(a)}${dot ? `<span class="ico-dot ${dot}"></span>` : ''}</div>
      <div class="app-tx"><div class="app-name">${esc(a.Name)}</div><div class="app-desc" title="${esc(L(a, 'Desc'))}">${esc(L(a, 'Desc'))}</div>
        ${!a.Installed && (a.Note || '').startsWith('bundled:') ? `<div class="app-note">${ic('check')}${esc(t('java_bundled', a.Note.slice(8)))}</div>` : ''}</div>
      <div class="check ${sel ? 'on' : ''}" data-act="sel" data-id="${esc(a.ID)}" title="${t('act_select')}">${ic('check')}</div></div>
    <div class="app-bot">${status}${btn}</div>
    <div class="prog"><div class="txt"><b>${esc(stepText(p))}</b><span>${pct !== null ? pct + '%' : ''}</span></div>
      <div class="bar ${pct === null ? 'indet' : ''}"><i style="width:${pct ?? 0}%"></i></div></div>
  </div>`;
}

// Apps die je mag zien: zoeken + weergave (alle / geïnstalleerd / updates) + alleen de tool van je eigen merk.
function visibleApps() {
  const q = S.q.trim().toLowerCase();
  const brand = S.hw && S.hw.OEM ? S.hw.OEM.Brand : '';
  return S.apps.filter(a => {
    if (a.Cat === 'oem' && a.Icon !== brand) return false;
    if (q && !(a.Name.toLowerCase().includes(q) || t('cat_' + a.Cat).toLowerCase().includes(q) || L(a, 'Desc').toLowerCase().includes(q))) return false;
    if (S.view === 'installed') return a.Installed;
    if (S.view === 'updates') return !!a.Available;
    return true;
  });
}

const APP_CATS = ['oem', 'browser', 'chat', 'ai', 'games', 'gear', 'media', 'tools', 'system', 'network', 'dev', 'runtime'];
const CAT_ICON = { updates: 'up', ai: 'spark', popular: 'spark', oem: 'shield', browser: 'globe', chat: 'apps', games: 'game', gear: 'zap', media: 'play', tools: 'tweaks', system: 'cpu', network: 'globe', dev: 'apps', runtime: 'pkg' };

function renderApps() {
  const pg = $('#p-apps');
  if (pg.dataset.built !== LANG) {
    pg.dataset.built = LANG;
    pg.innerHTML = `
      <div class="head"><div><h1>${t('apps_title')}</h1><p>${t('apps_sub')}</p></div>
        <div class="actions">
          <label class="search">${ic('search')}<input id="q" placeholder="${t('search_app')}" spellcheck="false" value="${esc(S.q)}">
            <button class="search-x ${S.q ? 'on' : ''}" id="qx" type="button" title="${t('clear_search')}">${ic('x')}</button></label>
          <button class="btn icon" id="rescan" title="${t('rescan')}">${ic('refresh')}</button>
        </div></div>
      <div id="appBanner"></div>
      <div class="views" id="views"></div>
      <div class="cat-bar" id="catBar"></div>
      <div id="appSections"></div>`;
    const clearQ = () => { S.q = ''; $('#q').value = ''; $('#qx').classList.remove('on'); renderAppGrid(); $('#q').focus(); };
    $('#q').addEventListener('input', e => { S.q = e.target.value; $('#qx').classList.toggle('on', !!S.q); renderAppGrid(); });
    $('#q').addEventListener('keydown', e => { if (e.key === 'Escape' && S.q) { e.stopPropagation(); clearQ(); } });
    $('#qx').addEventListener('click', e => { e.preventDefault(); clearQ(); });
    $('#rescan').addEventListener('click', scanApps);
    pg.addEventListener('scroll', spyCats, { passive: true });
  }
  renderAppGrid();
}

function renderAppGrid() {
  const list = visibleApps();
  const all = S.apps.filter(a => a.Cat !== 'oem' || (S.hw && S.hw.OEM && a.Icon === S.hw.OEM.Brand));
  const nInst = all.filter(a => a.Installed).length, nUpd = all.filter(a => a.Available).length;
  $('#rescan').classList.toggle('spin', S.scanning);
  // weergave-knoppen
  $('#views').innerHTML = [['all', t('view_all'), all.length], ['installed', t('view_installed'), nInst], ['updates', t('view_updates'), nUpd], ['pc', t('view_pc'), S.progs ? S.progs.length : '·']]
    .map(([k, n, c]) => `<button class="view ${S.view === k ? 'on' : ''} ${k === 'updates' && c ? 'has' : ''}" data-view="${k}">${n}<span>${S.appsLoaded ? c : '…'}</span></button>`).join('')
    + `<span class="views-tip">${t('apps_tip')}</span>`;
  if (S.view === 'pc') { $('#appBanner').innerHTML = ''; $('#catBar').style.display = 'none'; return renderPrograms(); }
  // bovenaan: winget ontbreekt, of alle updates in één keer
  let banner = '';
  if (S.sys && S.sys.Winget === false) banner = `<div class="banner warn"><div class="ic">${ic('alert')}</div><div class="tx"><b>${t('winget_missing')}</b><span>${t('winget_missing_sub')}</span></div>
    ${S.wg ? `<span class="pill pink">${S.wg.phase === 'fout' ? t('step_failed') : esc(t('step_' + S.wg.step))}</span>` : `<button class="btn primary sm" data-act="wgInstall">${ic('down')}${t('wg_install')}</button>`}</div>`;
  else if (nUpd && S.view !== 'installed') banner = `<div class="banner pink upd-banner"><div class="ic">${ic('up')}</div><div class="tx"><b>${t('upd_banner', nUpd)}</b><span>${t('upd_banner_sub')}</span></div>
    <button class="btn primary sm" data-act="updAll" ${S.busy ? 'disabled' : ''}>${ic('down')}${t('update_all_n', nUpd)}</button></div>`;
  $('#appBanner').innerHTML = banner;
  // secties: Populair (alleen bij "alle apps" zonder zoeken), dan per categorie
  const secs = [];
  if (S.view === 'all' && !S.q.trim()) secs.push(['popular', list.filter(a => a.Pop > 0).sort((x, y) => x.Pop - y.Pop).slice(0, 15)]);
  if (S.view === 'updates') { if (list.length) secs.push(['updates', list]); } // updates: één lijst
  else for (const c of APP_CATS) {
    const items = list.filter(a => a.Cat === c).sort((x, y) => (x.Pop || 999) - (y.Pop || 999));
    if (items.length) secs.push([c, items]);
  }
  $('#catBar').innerHTML = secs.length > 1 ? secs.map(([c]) => `<button class="chip" data-jump="${c}">${t(c === 'popular' ? 'cat_popular' : 'cat_' + c)}</button>`).join('') : '';
  $('#catBar').style.display = secs.length > 1 ? '' : 'none';
  $('#appSections').innerHTML = secs.length ? secs.map(([c, items]) => `
    <section class="app-sec" id="sec-${c}">
      <h2>${ic(CAT_ICON[c] || 'apps')}${t(c === 'popular' ? 'cat_popular' : c === 'updates' ? 'view_updates' : 'cat_' + c)}<span>${items.length}</span>${c === 'popular' ? `<em>${t('popular_sub')}</em>` : ''}</h2>
      <div class="grid ${S.sel.size ? 'selecting' : ''}">${items.map(appCard).join('')}</div>
    </section>`).join('')
    : `<div class="empty-state">${MOON}<h3>${S.view === 'updates' ? t('no_updates') : S.view === 'installed' && !S.q ? t('none_installed') : t('no_match')}</h3>${S.view === 'updates' ? t('no_updates_sub') : t('no_match_sub')}</div>`;
  spyCats();
}

// ---------- alle programma's op deze pc (ook niet in Dusk) ----------
async function loadPrograms() {
  try { S.progs = await API().Programs() || []; } catch { S.progs = []; }
  if (S.page === 'apps') renderAppGrid();
}

const fmtYMD = d => /^\d{8}$/.test(d || '') ? fmtDate(`${d.slice(0, 4)}-${d.slice(4, 6)}-${d.slice(6, 8)}`) : '';

function renderPrograms() {
  const box = $('#appSections');
  if (!S.progs) { box.innerHTML = `<div class="empty-state">${MOON}<h3>${t('scanning')}</h3></div>`; if (!S.progLoading) { S.progLoading = true; loadPrograms().then(() => { S.progLoading = false; }); } return; }
  const q = S.q.trim().toLowerCase();
  let list = S.progs.filter(p => !q || (p.Name + ' ' + p.Publisher).toLowerCase().includes(q));
  const sorters = { name: (a, b) => a.Name.localeCompare(b.Name), size: (a, b) => b.SizeMB - a.SizeMB, date: (a, b) => (b.InstallDate || '').localeCompare(a.InstallDate || '') };
  list = [...list].sort((a, b) => (a.Protected - b.Protected) || sorters[S.progSort](a, b));
  box.innerHTML = `<div class="prog-head"><span>${t('pc_sub')}</span>
      <div class="sorts">${['name', 'size', 'date'].map(k => `<button class="chip ${S.progSort === k ? 'on' : ''}" data-psort="${k}">${t('sort_' + k)}</button>`).join('')}</div></div>
    <div class="card prog-list">${list.length ? list.map(p => {
      const pr = S.progProg[p.Key], sel = S.progSel.has(p.Key);
      const icon = p.Icon && !/[\\:]/.test(p.Icon) ? `<div class="ico" style="--c:${(S.apps.find(a => a.ID === p.CatalogID) || {}).Color || '#5b3d8f'}"><i style="--m:url('icons/${esc(p.Icon)}.svg')"></i></div>`
        : `<div class="ico" style="--c:#4a3a78">${esc(p.Name.replace(/[^A-Za-z0-9]/g, '').slice(0, 2))}</div>`;
      return `<div class="prow ${sel ? 'sel' : ''} ${p.Protected ? 'prot' : ''}" ${p.Protected ? `title="${t('prot_tip')}"` : `data-act="psel" data-id="${esc(p.Key)}"`}>
        ${p.Protected ? `<div class="lock">${ic('lock')}</div>` : `<div class="check ${sel ? 'on' : ''}">${ic('check')}</div>`}
        ${icon}<div class="tx"><b>${esc(p.Name)}</b><span>${esc([p.Publisher, p.Version].filter(Boolean).join(' · '))}</span></div>
        <span class="psize">${p.SizeMB ? fmtSize(p.SizeMB * 1048576) : ''}</span><span class="pdate">${fmtYMD(p.InstallDate)}</span>
        ${pr ? `<span class="pill ${pr.phase === 'klaar' ? 'ok' : pr.phase === 'fout' ? 'bad' : 'pink'}">${stepText(pr)}</span>` : ''}</div>`;
    }).join('') : `<div class="drv-state">${t('no_match')}</div>`}</div>`;
}

// ---------- verwijderen met opruimen ----------
async function uninstallFlow(kind, ids) {
  const names = kind === 'apps' ? ids.map(id => (S.apps.find(a => a.ID === id) || {}).Name) : ids.map(k => (S.progs.find(p => p.Key === k) || {}).Name);
  const games = kind === 'apps' && ids.some(id => (S.apps.find(a => a.ID === id) || {}).Cat === 'games');
  const m = $('#modal');
  m.innerHTML = `<div class="dialog un-dlg"><h3>${t('un_title', ids.length)}</h3><p>${t('un_sub')}</p>
    <div class="un-apps">${names.map(n => `<span class="pill pink">${esc(n)}</span>`).join('')}</div>
    <div class="un-list"><div class="drv-state"><span class="spinner"></span>${t('un_checking')}</div></div>
    <p class="sd-note">${t('un_note')}</p>
    <div class="acts"><button class="btn ghost" data-v="0">${t('cancel')}</button><button class="btn danger" data-v="1" disabled>${ic('trash')}${t('uninstall_btn')}</button></div></div>`;
  m.classList.add('on');
  let plan = {};
  try { plan = (kind === 'apps' ? await API().UninstallPlan(ids) : await API().ProgramPlan(ids)) || {}; } catch { }
  const items = [];
  for (const [id, list] of Object.entries(plan)) for (const l of list || []) items.push({ ...l, id, on: l.Kind !== 'documents' && !(games && l.Kind === 'appdata') });
  const list = m.querySelector('.un-list');
  const draw = () => {
    const total = items.filter(x => x.on).reduce((a, x) => a + x.Size, 0);
    list.innerHTML = items.length ? `<div class="un-total">${t('un_total', fmtSize(total))}</div>` + items.map((x, i) => `<label class="un-item"><input type="checkbox" data-i="${i}" ${x.on ? 'checked' : ''}>
      <div class="tx"><b title="${esc(x.Path)}">${esc(x.Path)}</b><span>${t('kind_' + x.Kind)}</span></div><em>${x.Size ? fmtSize(x.Size) : ''}</em></label>`).join('')
      : `<div class="drv-state ok">${ic('check')}${t('un_none')}</div>`;
  };
  draw();
  m.querySelector('[data-v="1"]').disabled = false;
  list.onchange = e => { const i = e.target.dataset.i; if (i !== undefined) { items[i].on = e.target.checked; draw(); } };
  const go = await new Promise(res => { m.onclick = e => { const b = e.target.closest('[data-v]'); if (b) res(b.dataset.v === '1'); else if (e.target === m) res(false); }; });
  m.classList.remove('on'); m.onclick = null;
  if (!go) return;
  const paths = items.filter(x => x.on).map(x => x.Path);
  const msg = kind === 'apps' ? await API().UninstallApps(ids, paths) : await API().UninstallPrograms(ids, paths);
  if (msg) return toast(t(msg));
  setBusy(true);
  S.sel.clear(); S.progSel.clear();
  render();
}

function showUninstallResult(d) {
  const m = $('#modal');
  const removed = d.removed || [];
  m.innerHTML = `<div class="dialog un-dlg"><h3>${t('un_done_title')}</h3>
    <p>${t('un_done', d.ok, fmtSize(d.freed || 0))}${d.fail ? '<br>' + t('un_failed', d.fail) : ''}</p>
    ${removed.length ? `<div class="sd-note">${t('un_removed_list')}</div><div class="un-list done">${removed.map(p => `<div class="un-item"><span class="okc">${ic('check')}</span><div class="tx"><b>${esc(p)}</b></div></div>`).join('')}</div>` : ''}
    <div class="acts"><button class="btn primary" data-v="1">${t('ok_btn')}</button></div></div>`;
  m.classList.add('on');
  m.onclick = e => { if (e.target.closest('[data-v]') || e.target === m) { m.classList.remove('on'); m.onclick = null; } };
}

// Welke categorie staat nu bovenaan? Die licht op in de balk.
function spyCats() {
  const pg = $('#p-apps'), bar = $('#catBar');
  if (!pg || !bar || !bar.children.length) return;
  const top = pg.getBoundingClientRect().top + bar.offsetHeight + 40;
  let cur = null;
  for (const sec of pg.querySelectorAll('.app-sec')) if (sec.getBoundingClientRect().top <= top) cur = sec.id.slice(4);
  if (!cur) cur = (pg.querySelector('.app-sec') || { id: 'sec-' }).id.slice(4);
  bar.querySelectorAll('.chip').forEach(ch => ch.classList.toggle('on', ch.dataset.jump === cur));
}

function updateCard(id, flash) {
  const a = S.apps.find(x => x.ID === id);
  if (!a) return;
  document.querySelectorAll(`.app[data-id="${CSS.escape(id)}"]`).forEach(el => {
    const tmp = document.createElement('div');
    tmp.innerHTML = appCard(a);
    const n = tmp.firstElementChild;
    if (flash) n.classList.add(flash);
    el.replaceWith(n);
  });
}

async function scanApps() {
  if (S.scanning) return;
  S.lastRefresh = Date.now();
  S.scanning = true;
  if (S.page === 'apps') renderAppGrid();
  try {
    S.apps = await API().ScanApps();
    S.appsErr = null;
  } catch (e) {
    S.appsErr = String(e);
    toast(t('scan_failed', e), 'bad');
  }
  S.appsLoaded = true;
  S.scanning = false;
  renderNav();
  render();
}

// Voor een update: kijk of andere apps de bestanden gebruiken (bv. Discord met de
// virtuele camera van OBS). Zo ja, vraag of Dusk ze mag sluiten.
async function checkBlockers(ids, retry) {
  const inst = ids.filter(id => (S.apps.find(a => a.ID === id) || {}).Installed);
  if (!inst.length) return true;
  inst.forEach(id => { S.prog[id] = { id, phase: 'bezig', pct: -1, step: 'checking' }; updateCard(id); });
  let map = {};
  try { map = await API().InUse(inst) || {}; } catch { }
  inst.forEach(id => { delete S.prog[id]; updateCard(id); });
  const blockers = {}, hit = [];
  for (const [id, list] of Object.entries(map)) {
    if (!list || !list.length) continue;
    hit.push((S.apps.find(a => a.ID === id) || { Name: id }).Name);
    for (const b of list) blockers[b.Name] = b;
  }
  const bl = Object.values(blockers);
  if (!bl.length) return !retry;
  const names = hit.length > 1 ? hit.slice(0, -1).join(', ') + ' ' + t('and') + ' ' + hit[hit.length - 1] : hit[0];
  const choice = await blockerDialog(bl, names, retry);
  if (choice === 'close') {
    toast(t('closing_apps'));
    await API().CloseApps(bl.map(b => b.Name));
    return true;
  }
  return choice === 'anyway';
}

function blockerDialog(list, names, retry) {
  return new Promise(res => {
    const m = $('#modal');
    m.innerHTML = `<div class="dialog blk-dlg"><h3>${retry ? t('blk_retry_title') : t('blk_title')}</h3>
      <p>${esc(retry ? t('blk_retry_text', names) : t('blk_text', names))}</p>
      <div class="blk-list">${list.map(b => `<div class="blk"><span class="blk-ic">${esc(b.Desc.slice(0, 1).toUpperCase())}</span><b>${esc(b.Desc)}</b><em>${esc(b.Name)}.exe</em></div>`).join('')}</div>
      <p class="blk-note">${t('blk_note')}</p>
      <div class="acts"><button class="btn ghost" data-v="cancel">${t('cancel')}</button>
        ${retry ? '' : `<button class="btn" data-v="anyway">${t('blk_anyway')}</button>`}
        <button class="btn primary" data-v="close">${retry ? t('blk_retry') : t('blk_close')}</button></div></div>`;
    m.classList.add('on');
    const close = v => { m.classList.remove('on'); m.onclick = null; document.removeEventListener('keydown', key); res(v); };
    const key = e => { if (e.key === 'Escape') close('cancel'); };
    document.addEventListener('keydown', key);
    m.onclick = e => { const b = e.target.closest('[data-v]'); if (b) close(b.dataset.v); else if (e.target === m) close('cancel'); };
  });
}

async function install(ids, retry) {
  if (!ids.length) return;
  if (S.busy) return toast(t('busy'));
  if (S.sys && S.sys.Winget === false) { if (S.page !== 'apps') goPage('apps'); return toast(t('winget_missing_toast'), 'bad'); }
  if (!(await checkBlockers(ids, retry))) return;
  const msg = await API().InstallApps(ids);
  if (msg) return toast(t(msg));
  setBusy(true);
  S.sel.clear();
  render();
}

function renderSel() {
  const bar = $('#selbar');
  if (S.page === 'apps' && S.view === 'pc' && S.progSel.size) {
    bar.innerHTML = `<b>${t('n_selected', S.progSel.size)}</b>
      <button class="btn danger sm" data-act="progRemove" ${S.busy ? 'disabled' : ''}>${ic('trash')}${t('uninstall_btn')}</button>
      <button class="btn ghost sm icon" data-act="clearProg" title="${t('clear_sel')}">${ic('x')}</button>`;
    bar.classList.add('on');
  } else if (S.page === 'apps' && S.sel.size) {
    const list = [...S.sel].map(id => S.apps.find(a => a.ID === id)).filter(Boolean);
    const todo = list.filter(a => !a.Installed || a.Available).length;
    const rem = list.filter(a => a.Installed && a.Cat !== 'runtime').length;
    bar.innerHTML = `<b>${t('n_selected', S.sel.size)}</b><span style="color:var(--muted)">${t('n_todo', todo)}</span>
      <button class="btn primary sm" data-act="installSel" ${todo && !S.busy ? '' : 'disabled'}>${ic('down')}${t('install_update')}</button>
      ${rem ? `<button class="btn danger sm" data-act="appRemove" ${S.busy ? 'disabled' : ''}>${ic('trash')}${t('uninstall_btn')} (${rem})</button>` : ''}
      <button class="btn ghost sm icon" data-act="clearSel" title="${t('clear_sel')}">${ic('x')}</button>`;
    bar.classList.add('on');
  } else if (S.page === 'clean' && S.bloatSel.size) {
    bar.innerHTML = `<b>${t('n_selected', S.bloatSel.size)}</b>
      <button class="btn danger sm" data-act="removeSel" ${S.busy ? 'disabled' : ''}>${ic('trash')}${t('remove')}</button>
      <button class="btn ghost sm icon" data-act="clearBloat" title="${t('clear_sel')}">${ic('x')}</button>`;
    bar.classList.add('on');
  } else if (S.page === 'hardware' && (S.drvSel.size || S.oemSel.size)) {
    bar.innerHTML = `<b>${t('n_selected', S.drvSel.size + S.oemSel.size)}</b><span style="color:var(--muted)">${S.sys && S.sys.Admin ? '' : t('drv_admin_note')}</span>
      <button class="btn primary sm" data-act="drvInstall" ${S.busy ? 'disabled' : ''}>${ic('down')}${t('install_selected')}</button>
      <button class="btn ghost sm icon" data-act="clearDrv" title="${t('clear_sel')}">${ic('x')}</button>`;
    bar.classList.add('on');
  } else bar.classList.remove('on');
}

// ---------- tweaks ----------
function renderTweaks() {
  const admin = S.sys && S.sys.Admin;
  const groups = GROUP_ORDER.filter(g => S.tweaks.some(x => x.Group === g));
  const onN = S.tweaks.filter(x => x.On).length;
  const hasAdm = S.tweaks.some(x => x.Admin);
  // vaste verdeling over twee kolommen (op aantal tweaks), zodat uitklappen niks laat verspringen
  const cols = [[], []], size = [0, 0];
  for (const g of groups) { const c = size[0] <= size[1] ? 0 : 1; cols[c].push(g); size[c] += S.tweaks.filter(x => x.Group === g).length + 2; }
  const groupHtml = g => `
      <div class="card tw-group"><h3>${ic(GROUP_ICON[g] || 'spark')}${t('grp_' + g)}</h3>
        ${S.tweaks.filter(x => x.Group === g).map(x => {
          const locked = x.Admin && !admin, wait = S.tweakWait.has(x.ID), note = L(x, 'Note');
          const open = S.twOpen.has(x.ID);
          return `<div class="tw ${open ? 'open' : ''} ${locked ? 'locked' : ''}" data-tw="${x.ID}"><div class="tw-row"><div class="tx"><div class="nm">${esc(L(x, 'Name'))}${x.Admin ? `<span class="pill adm">${t('admin')}</span>` : ''}
              <button class="info-btn" data-info="${x.ID}" title="${t('more_info')}">i</button></div>
            <div class="ds">${esc(L(x, 'Desc'))}</div>${note ? `<div class="nt">${esc(note)}</div>` : ''}</div>
            <label class="switch ${locked ? 'locked' : ''} ${wait ? 'wait' : ''}" title="${locked ? t('start_admin_tip') : ''}">
              <input type="checkbox" data-tweak="${x.ID}" ${x.On ? 'checked' : ''} ${locked || wait ? 'disabled' : ''}><span></span></label></div>
            ${open ? `<div class="tw-info"><b>${t('info_what')}</b><p>${esc(L(x, 'Info'))}</p>
              <b>${t('info_changes')}</b><div class="regs">${(x.Changes || []).map(c => `<code>${esc(c)}</code>`).join('')}</div></div>` : ''}</div>`;
        }).join('')}
      </div>`;
  const pg = $('#p-tweaks'), keep = pg.scrollTop;
  $('#p-tweaks').innerHTML = `
    <div class="head"><div><h1>${t('tweaks_title')}</h1><p>${t('tweaks_sub')}</p></div>
      <div class="actions"><span class="pill pink">${t('n_of_on', onN, S.tweaks.length)}</span>
        <button class="btn sm" data-act="restore" title="${t('restore_sub')}">${ic('shield')}${t('restore_btn')}</button></div></div>
    ${S.explorerPending ? `<div class="banner pink"><div class="ic">${ic('folder')}</div><div class="tx"><b>${t('restart_explorer')}</b><span>${t('restart_explorer_sub')}</span></div><button class="btn primary sm" data-act="explorer">${t('restart_now')}</button></div>` : ''}
    ${!admin && hasAdm ? `<div class="banner warn"><div class="ic">${ic('lock')}</div><div class="tx"><b>${t('admin_needed')}</b><span>${t('admin_needed_sub')}</span></div><button class="btn sm" data-act="admin">${ic('shield')}${t('run_admin')}</button></div>` : ''}
    <div class="tw-grid">${[0, 1].map(c => `<div class="tw-col">${cols[c].map(groupHtml).join('')}</div>`).join('')}</div>`;
  pg.scrollTop = keep;
}

async function toggleTweak(id, want) {
  S.tweakWait.add(id);
  renderTweaks();
  try {
    const x = await API().SetTweak(id, want);
    const i = S.tweaks.findIndex(y => y.ID === id);
    S.tweaks[i] = x;
    if (x.On !== want) toast(t('tweak_partly', L(x, 'Name')), 'bad');
    else {
      toast(t('tweak_is', L(x, 'Name'), want), 'ok');
      if (x.Explorer) S.explorerPending = true;
    }
  } catch (e) {
    toast(t('error', e), 'bad');
  }
  S.tweakWait.delete(id);
  renderNav();
  if (S.page === 'tweaks') renderTweaks();
}

// ---------- hardware ----------
async function loadHardware() {
  S.lastRefresh = Date.now();
  try { S.hw = await API().Hardware(); S.hwErr = null; } catch (e) { S.hwErr = String(e); }
  renderNav();
  if (S.page === 'hardware') renderHardware();
  if (S.page === 'apps') renderAppGrid();
}

// Alle drivers controleren: Windows Update + NVIDIA/AMD + HP-lijst (±30 s).
async function checkDrivers() {
  if (S.drvChecking || !S.hw) return;
  S.drvChecking = true;
  if (S.page === 'hardware') renderHardware();
  const amd = (S.hw.Software || []).find(x => x.Name === 'AMD Software');
  const gpus = S.hw.GPUs.filter(g => !g.Virtual).map(g => ({ Name: g.Name, Vendor: g.Vendor, Ver: g.Ver, Software: g.Vendor === 'amd' && amd ? amd.Ver : '' }));
  try { S.drvStat = await API().CheckDrivers(gpus); } catch (e) { S.drvStat = { Updates: [], GPU: {}, OEM: [], Err: String(e) }; }
  S.drvChecking = false;
  S.drvProg = {}; S.oemProg = {};
  renderNav();
  if (S.page === 'hardware') { renderHardware(); renderSel(); }
  if (S.page === 'home') renderHome();
}

const daysSince = iso => iso ? Math.floor((Date.now() - new Date(iso + 'T00:00:00').getTime()) / 864e5) : null;
const fmtDate = iso => iso ? new Date(iso + 'T00:00:00').toLocaleDateString(locale(), { day: 'numeric', month: 'short', year: 'numeric' }) : '-';
const shortBoard = b => (b || '').replace(/Micro-Star International Co\.,? ?Ltd\.?/i, 'MSI').replace(/ASUSTeK COMPUTER INC\.?/i, 'ASUS')
  .replace(/Gigabyte Technology Co\.,? ?Ltd\.?/i, 'Gigabyte').replace(/ASRock\w*/i, 'ASRock').replace(/\s*\(MS-\w+\)/i, '').replace(/\s+/g, ' ').trim();
const fmtSize = b => b >= 1e9 ? (b / 1e9).toFixed(1) + ' GB' : b >= 1e6 ? Math.round(b / 1e6) + ' MB' : b > 0 ? Math.max(1, Math.round(b / 1e3)) + ' KB' : '0 MB';
const VENDOR = {
  amd: { color: '#ED1C24', url: 'https://www.amd.com/en/support/download/drivers.html', label: 'open_amd' },
  nvidia: { color: '#76B900', url: 'https://www.nvidia.com/en-us/software/nvidia-app/', label: 'open_nvidia' },
  intel: { color: '#0071C5', url: 'https://www.intel.com/content/www/us/en/support/intel-driver-support-assistant.html', label: 'open_intel' },
};
const BRAND = { hp: 'HP', dell: 'Dell', lenovo: 'Lenovo', asus: 'ASUS', acer: 'Acer', msi: 'MSI', microsoft: 'Microsoft' };

function hwCard(icon, label, value, extra = '') {
  return `<div class="card hw"><div class="hw-ic">${ic(icon)}</div><div class="hw-tx"><div class="k">${label}</div><div class="v" title="${esc(value)}">${esc(value || '…')}</div>${extra}</div></div>`;
}

function visibleProblems() {
  const ign = S.set.ignored || [];
  return S.hw ? (S.hw.Devices || []).filter(d => d.Problem > 0 && !ign.includes(d.Name.trim())) : [];
}

function boardUrl(h) {
  const maker = (h.BoardMaker || '').toLowerCase(), prod = (h.BoardProduct || '').replace(/\s*\(.*?\)\s*/g, '').trim();
  if (h.OEM && h.OEM.Brand !== 'custom' && h.OEM.SupportURL) return h.OEM.SupportURL;
  if (maker.includes('micro-star') || maker.includes('msi')) return `https://www.msi.com/Motherboard/${prod.replace(/\s+/g, '-')}/support`;
  if (maker.includes('gigabyte')) return `https://www.gigabyte.com/Search?kw=${encodeURIComponent(prod)}`;
  return `https://www.google.com/search?q=${encodeURIComponent(shortBoard(h.Board) + ' drivers support')}`;
}

function makerUrl(maker, h) {
  const m = (maker || '').toLowerCase();
  if (h.OEM && !['custom', 'other'].includes(h.OEM.Brand) && h.OEM.SupportURL) return h.OEM.SupportURL; // laptops: altijd de maker van de laptop
  if (m.includes('advanced micro') || /^amd\b/.test(m)) return VENDOR.amd.url;
  if (m.includes('intel')) return VENDOR.intel.url;
  if (m.includes('nvidia')) return VENDOR.nvidia.url;
  if (m.includes('logitech')) return 'https://support.logi.com/hc/en-us/sections/360004305013';
  return boardUrl(h);
}

// Status van één apparaat: missing, problem, update, latest, nonewer, old, windows, checking, unknown
function devStatus(d) {
  if (d.Problem === 28 || d.Basic) return { k: 'missing', wu: wuFor(d) };
  if (d.Problem > 0) return { k: 'problem', wu: wuFor(d) };
  if (S.drvChecking) return { k: 'checking' };
  if (!S.drvStat) return { k: 'unknown' };
  const wu = wuFor(d);
  if (wu) return { k: 'update', wu };
  const oem = (S.drvStat.OEM || []).find(u => u.Device === d.Name);
  if (oem) return { k: 'update', oem };
  if (d.Cat === 'gpu') {
    const g = (S.drvStat.GPU || {})[d.Name];
    if (g && g.Known) return g.Newer ? { k: 'update', url: g.URL, dl: g.Download, ver: g.Version } : { k: 'latest', ver: g.Version };
  }
  if (d.Generic) return { k: 'windows' };
  const days = daysSince(d.Date);
  return { k: days !== null && days > 3 * 365 ? 'old' : 'nonewer' };
}

function wuFor(d) {
  if (!S.drvStat || !S.drvStat.Updates) return null;
  const ids = (d.HWIDs || []).map(x => x.toLowerCase());
  return S.drvStat.Updates.find(u => u.HWID && ids.includes(u.HWID.toLowerCase())) || null;
}

// Bijwerken in Dusk: NVIDIA installeren, AMD Software openen, Intel-tool, of de tool van het merk.
function gpuAction(d, st) {
  const p = S.gpuProg[d.Name];
  if (p) return `<span class="pill ${p.phase === 'klaar' ? 'ok' : p.phase === 'fout' ? 'bad' : 'pink'}">${stepText(p)}${p.pct > 0 && p.pct < 100 ? ' ' + p.pct + '%' : ''}</span>`;
  if (st.wu) {
    const w = S.drvProg[st.wu.ID];
    if (w) return `<span class="pill ${w.phase === 'klaar' ? 'ok' : w.phase === 'fout' ? 'bad' : 'pink'}">${stepText(w)}</span>`;
    return `<button class="btn primary sm" data-act="wuOne" data-id="${esc(st.wu.ID)}">${ic('down')}${t('act_install')}</button>`;
  }
  if (st.dl) return `<button class="btn primary sm" data-act="gpuInstall" data-name="${esc(d.Name)}" data-dl="${esc(st.dl)}">${ic('down')}${t('act_install')}</button>`;
  if (/amd|advanced micro/i.test(d.Maker)) return `<button class="btn primary sm" data-act="amdOpen">${ic('ext')}${t('act_amd')}</button>`;
  return `<button class="btn primary sm" data-url="${esc(st.url)}">${ic('ext')}${t('act_download')}</button>`;
}

function keepAction(d, h, old) {
  const o = h.OEM || {}, cls = old ? '' : 'ghost';
  if (['dell', 'lenovo'].includes(o.Brand)) return `<button class="btn ${cls} sm" data-act="oemRun">${ic('refresh')}${t('act_oemrun', o.ToolName)}</button>`;
  if (/intel/i.test(d.Maker)) return `<button class="btn ${cls} sm" data-act="intelCheck">${ic('search')}${t('act_intel')}</button>`;
  if (/advanced micro|^amd/i.test(d.Maker) && d.Cat === 'gpu') return `<button class="btn ${cls} sm" data-act="amdOpen">${ic('ext')}${t('act_amd')}</button>`;
  return `<button class="btn ghost sm" data-url="${esc(makerUrl(d.Maker, h))}">${ic('ext')}${t('act_maker')}</button>`;
}

const STATUS_PILL = { missing: 'bad', problem: 'bad', update: 'upd', latest: 'ok', nonewer: 'ok', old: 'upd', windows: '', checking: 'pink', unknown: '' };

function devRow(d, h) {
  const st = devStatus(d), days = daysSince(d.Date);
  let act = '';
  if (st.wu) {
    const p = S.drvProg[st.wu.ID];
    act = p ? `<span class="pill ${p.phase === 'klaar' ? 'ok' : p.phase === 'fout' ? 'bad' : 'pink'}">${stepText(p)}</span>`
      : `<div class="check ${S.drvSel.has(st.wu.ID) ? 'on' : ''}" data-act="dsel" data-id="${esc(st.wu.ID)}" title="${t('act_select')}">${ic('check')}</div>`;
  } else if (st.oem) act = `<button class="btn sm" data-goto-oem>${ic('down')}${t('act_download')}</button>`;
  else if (st.k === 'update' && st.url) act = gpuAction(d, st);
  else if (st.k === 'missing' || st.k === 'problem') {
    const hw = (d.HWIDs || [])[0] || d.Name;
    act = `<button class="btn sm" data-url="https://www.google.com/search?q=${encodeURIComponent(hw + ' driver')}">${ic('search')}${t('act_search')}</button>
      <button class="btn ghost sm" data-act="ignore" data-name="${esc(d.Name.trim())}">${t('ignore')}</button>`;
  } else if (['old', 'nonewer', 'windows', 'unknown'].includes(st.k)) act = keepAction(d, h, st.k === 'old');
  const meta = d.Problem ? (I18N[LANG]['code_' + d.Problem] ? t('code_' + d.Problem) : t('code_x', d.Problem))
    : [d.Maker, d.Ver && `<code>${esc(d.Ver)}</code>`, d.Date && `${fmtDate(d.Date)} · ${t('days_ago', days)}`].filter(Boolean).join(' · ');
  const extra = st.wu ? ` · ${t('newest_is', esc(st.wu.Date ? fmtDate(st.wu.Date) : st.wu.Title))}` : st.oem ? ` · ${t('newest_is', esc(st.oem.Version))}` : st.ver && st.k === 'update' ? ` · ${t('newest_is', esc(st.ver))}` : '';
  return `<div class="dev"><div class="tx"><b title="${esc(d.Name)}">${esc(d.Name)}</b><span class="meta">${meta.replace(/</g, '<')}${extra}</span></div>
    <span class="pill ${STATUS_PILL[st.k]}"><i></i>${t('s_' + st.k)}</span><div class="dev-act">${act}</div></div>`;
}

function renderHardware() {
  const h = S.hw, s = S.sys || {};
  let body;
  if (S.hwErr) body = `<div class="banner warn"><div class="ic">${ic('alert')}</div><div class="tx"><b>${t('error', '')}</b><span>${esc(S.hwErr)}</span></div></div>`;
  else if (!h) body = `<div class="empty-state">${MOON}<h3>${t('scanning')}</h3></div>`;
  else {
    const o = h.OEM || {}, devs = h.Devices || [];
    const real = h.GPUs.filter(g => !g.Virtual);
    // jouw pc
    const brand = BRAND[o.Brand] || '';
    const mk = (o.Maker || '').trim(), md = (o.Model || '').trim();
    const title = o.Brand === 'custom' ? t('custom_pc') : (md.toLowerCase().startsWith(mk.toLowerCase()) ? md : `${mk} ${md}`.trim());
    const tool = o.Tool ? S.apps.find(a => a.ID === o.Tool) : null;
    const hasIntel = devs.some(d => /intel/i.test(d.Maker)) || /intel/i.test(h.CPU);
    const intelTool = S.apps.find(a => a.ID === 'Intel.IntelDriverAndSupportAssistant');
    const pcCard = `<div class="card pc-card">
      <div class="pc-ic">${ic(o.Laptop ? 'laptop' : 'desktop')}</div>
      <div class="pc-tx"><span class="k">${t('your_pc')} · ${o.Laptop ? t('laptop') : t('desktop')}</span><b>${esc(o.Brand === 'custom' ? shortBoard(h.Board) : title)}</b>
        ${o.Brand === 'custom' ? `<span>${t('custom_pc')}</span>` : ''}
        ${o.ToolName ? `<p>${esc(t('oem_tool', brand, o.ToolName))}</p>` : ''}</div>
      <div class="pc-acts">
        ${o.Tool ? (tool && tool.Installed ? (['dell', 'lenovo'].includes(o.Brand) ? `<button class="btn primary sm" data-act="oemRun">${ic('refresh')}${esc(t('oem_run_all', o.ToolName))}</button>` : `<span class="pill ok"><i></i>${esc(o.ToolName)} · ${t('tool_installed')}</span>`) : `<button class="btn primary sm" data-act="one" data-id="${esc(o.Tool)}">${ic('down')}${esc(t('install_tool', o.ToolName))}</button>`) : ''}
        ${hasIntel && !['hp', 'dell', 'lenovo'].includes(o.Brand) ? (intelTool && intelTool.Installed ? '' : `<button class="btn sm" data-act="one" data-id="Intel.IntelDriverAndSupportAssistant">${ic('down')}${t('install_intel')}</button>`) : ''}
        <button class="btn ghost sm" data-url="${esc(o.Brand === 'custom' ? boardUrl(h) : (o.SupportURL || boardUrl(h)))}">${ic('ext')}${t('oem_support')}</button>
      </div></div>`;
    // geheugen
    const ram = h.RAM || [];
    const ramSpeed = ram.length ? Math.min(...ram.map(r => r.Speed || 0)) : 0;
    const rated = ram.length ? Math.max(...ram.map(r => r.Rated || 0)) : 0;
    const slow = rated && ramSpeed && rated > ramSpeed + 150;
    const sizes = [...new Set(ram.map(r => r.GB))];
    const ramText = ram.length ? `${s.RAM || ram.reduce((a, r) => a + r.GB, 0)} GB · ${sizes.length === 1 ? t('ram_sticks', ram.length, sizes[0]) : ram.map(r => r.GB + ' GB').join(' + ')}` : (s.RAM ? s.RAM + ' GB' : '');
    const ramExtra = ramSpeed ? `<div class="hw-sub ${slow ? 'warn' : 'ok'}">${t('ram_runs', ramSpeed)}${slow ? '' : rated ? ' · ' + t('ram_ok') : ''}</div>` : '';
    const bat = (h.Batteries || [])[0];
    const batPct = bat && bat.Design ? Math.round(bat.Full / bat.Design * 100) : 0;
    // videokaarten
    const gpus = real.map(g => {
      const d = devs.find(x => x.Name === g.Name) || { Name: g.Name, Cat: 'gpu', Ver: g.Ver, Date: g.Date, HWIDs: [] };
      const st = devStatus(d), v = VENDOR[g.Vendor], days = daysSince(g.Date);
      const note = g.Vendor === 'amd' ? t('gpu_note_amd') : g.Vendor === 'nvidia' ? t('gpu_note_nvidia') : g.Vendor === 'intel' ? t('gpu_note_intel') : t('gpu_tool_hint');
      return `<div class="card gpu">
        <div class="gpu-top"><div class="hw-ic" style="${v ? `background:${v.color}22;color:${v.color}` : ''}">${ic('gpu')}</div>
          <div class="hw-tx"><div class="k">${t('graphics')}</div><div class="v">${esc(g.Name.replace(/\s+/g, ' '))}</div></div>
          <span class="pill ${STATUS_PILL[st.k]}"><i></i>${t('s_' + st.k)}${st.ver && st.k === 'update' ? ' · ' + esc(st.ver) : ''}</span></div>
        <div class="gpu-rows">
          <div><span>${t('driver')}</span><b>${esc(g.Ver)}</b></div>
          <div><span>${t('driver_date')}</span><b>${fmtDate(g.Date)}${days !== null ? ` <em>· ${t('days_ago', days)}</em>` : ''}</b></div>
        </div>
        ${v ? `<div class="gpu-foot"><span>${note}</span>${st.k === 'update' || st.k === 'missing' && st.wu ? gpuAction(d, st) : keepAction(d, h, false)}</div>` : ''}
      </div>`;
    }).join('');
    // updates van de maker (HP)
    let oemBox = '';
    if (o.Deep) {
      let inner;
      if (S.drvChecking) inner = `<div class="drv-state"><span class="spinner"></span>${t('oem_checking')}</div>`;
      else if (!S.drvStat) inner = '';
      else if (S.drvStat.OEMErr) inner = `<div class="drv-state bad">${t('oem_err')} <code>${esc(S.drvStat.OEMErr)}</code></div>`;
      else if (!(S.drvStat.OEM || []).length) inner = `<div class="drv-state ok">${ic('check')}${t('oem_none')}</div>`;
      else inner = S.drvStat.OEM.map(u => {
        const p = S.oemProg[u.ID];
        return `<div class="drv ${S.oemSel.has(u.ID) ? 'sel' : ''}" data-act="osel" data-id="${esc(u.ID)}">
          <div class="check ${S.oemSel.has(u.ID) ? 'on' : ''}">${ic('check')}</div>
          <div class="tx"><b>${esc(u.Name)}</b><span>${esc(u.Category)} · ${esc(u.Installed)} → <b class="nv">${esc(u.Version)}</b> · ${fmtDate(u.Date)} · ${fmtSize(u.Size)}</span></div>
          ${p ? `<span class="pill ${p.phase === 'klaar' ? 'ok' : p.phase === 'fout' ? 'bad' : 'pink'}">${stepText(p)}${p.pct > 0 && p.pct < 100 ? ' ' + p.pct + '%' : ''}</span>` : ''}</div>`;
      }).join('');
      oemBox = `<div class="section-title" id="oemBox">${t('oem_updates')}</div>
        <div class="card drvbox"><div class="drv-head"><div><b>${t('oem_updates')}</b><span>${esc(t('oem_updates_sub', o.Model || 'HP'))}</span></div></div>${inner}</div>`;
    }
    // alle drivers
    const cats = ['gpu', 'chipset', 'audio', 'network', 'wifi', 'bluetooth', 'storage', 'usb', 'input', 'camera', 'other'];
    const ign = S.set.ignored || [];
    const shown = devs.filter(d => !(d.Problem > 0 && ign.includes(d.Name.trim())) && d.Cat !== 'gpu');
    const counts = { u: 0, m: 0, o: 0 };
    devs.forEach(d => { if (ign.includes(d.Name.trim())) return; const k = devStatus(d).k; if (k === 'update') counts.u++; if (k === 'missing' || k === 'problem') counts.m++; if (k === 'old') counts.o++; });
    const sum = S.drvChecking ? `<span class="spinner"></span>${t('dev_checking')}` : !S.drvStat ? '' : (counts.u + counts.m + counts.o ? esc(t('sum', counts.u, counts.m, counts.o)) : `${ic('check')}${t('sum_ok')}`);
    const groups = cats.map(c => {
      const list = shown.filter(d => d.Cat === c);
      if (!list.length) return '';
      return `<div class="dev-group"><h4>${t('cat_' + c)}</h4>${list.map(d => devRow(d, h)).join('')}</div>`;
    }).join('');
    const matched = new Set(devs.map(d => wuFor(d)).filter(Boolean).map(u => u.ID));
    const otherWU = S.drvStat ? (S.drvStat.Updates || []).filter(u => !matched.has(u.ID)) : [];
    const otherBox = otherWU.length ? `<div class="dev-group"><h4>${t('other_wu')}</h4>${otherWU.map(u => {
      const p = S.drvProg[u.ID];
      return `<div class="dev"><div class="tx"><b>${esc(u.Title)}</b><span class="meta">${esc([u.Maker, u.Class, fmtDate(u.Date)].filter(Boolean).join(' · '))}</span></div>
        <span class="pill upd"><i></i>${t('s_update')}</span><div class="dev-act">${p ? `<span class="pill ${p.phase === 'klaar' ? 'ok' : p.phase === 'fout' ? 'bad' : 'pink'}">${stepText(p)}</span>`
        : `<div class="check ${S.drvSel.has(u.ID) ? 'on' : ''}" data-act="dsel" data-id="${esc(u.ID)}">${ic('check')}</div>`}</div></div>`;
    }).join('')}</div>` : '';
    const vp = visibleProblems(), ignN = devs.filter(d => d.Problem > 0).length - vp.length;
    const disks = (h.Disks || []).map(d => `<div class="disk"><span>${esc(d.Name)}</span><em>${d.Size >= 1000 ? (d.Size / 1024).toFixed(1) + ' TB' : d.Size + ' GB'} ${esc(d.Type || '')}</em>
      <span class="pill ${d.Health === 'Healthy' ? 'ok' : 'bad'}"><i></i>${d.Health === 'Healthy' ? t('healthy') : t('unhealthy')}</span></div>`).join('');
    body = `
      ${slow ? `<div class="banner warn"><div class="ic">${ic('ram')}</div><div class="tx"><b>${t('ram_slow_title')}</b><span>${t('ram_slow_text', rated, ramSpeed)}</span></div></div>` : ''}
      ${pcCard}
      <div class="hw-grid">
        ${hwCard('cpu', t('processor'), shortCPU(h.CPU))}
        ${hwCard('board', t('motherboard'), shortBoard(h.Board), h.BIOS ? `<div class="hw-sub">BIOS ${esc(h.BIOS)}${h.BIOSDate ? ' · ' + fmtDate(h.BIOSDate) : ''}</div>` : '')}
        ${hwCard('ram', t('memory'), ramText, ramExtra)}
        ${(h.Monitors || []).length ? hwCard('desktop', t('screens'), h.Monitors.map(m => m.Name || m.Maker).filter(Boolean).join(', ')) : ''}
        ${bat ? hwCard('battery', t('battery'), bat.Charge + '%', batPct ? `<div class="hw-sub ${batPct < 70 ? 'warn' : 'ok'}">${t('battery_health', batPct)}</div>` : '') : ''}
      </div>
      ${gpus}
      ${oemBox}
      <div class="section-title">${t('dev_title')}</div>
      <div class="card drvbox"><div class="drv-head"><div><b>${t('dev_title')}</b><span>${t('dev_sub')} ${t('drv_wu_note')}</span><div class="dev-sum">${sum}</div></div>
        <button class="btn ${S.drvStat ? '' : 'primary'} sm ${S.drvChecking ? 'spin' : ''}" data-act="drvCheck" ${S.drvChecking ? 'disabled' : ''}>${ic(S.drvStat ? 'refresh' : 'search')}${t('dev_check')}</button></div>
        ${groups}${otherBox}
        ${vp.length ? `<p class="hint">${t('no_driver_hint')}</p>` : ''}
        <div class="drv-foot">${vp.length ? `<button class="btn sm" data-act="devmgr">${ic('ext')}${t('open_devmgr')}</button>` : ''}
          ${ignN ? `<span class="ign">${t('ignored_n', ignN)} · <a data-act="showIgnored">${t('show_ignored')}</a></span>` : ''}</div></div>
      ${disks ? `<div class="section-title">${t('drives')}</div><div class="card disks">${disks}</div>` : ''}`;
  }
  $('#p-hardware').innerHTML = `<div class="head"><div><h1>${t('hw_title')}</h1><p>${t('hw_sub')}</p></div>
    <div class="actions"><button class="btn icon" data-act="hwRefresh" title="${t('rescan')}">${ic('refresh')}</button></div></div>${body}`;
}

// ---------- opstart-apps ----------
async function loadStartup() {
  try { S.st = await API().Startup() || []; } catch (e) { S.st = []; }
  if (S.page === 'startup') renderStartup();
}

function renderStartup() {
  const admin = S.sys && S.sys.Admin;
  // alleen apps die veilig uit mogen; drivers en Windows-onderdelen blijven verborgen
  const list = S.st ? S.st.filter(x => !x.Essential) : null;
  const hiddenN = S.st ? S.st.length - list.length : 0;
  let body;
  if (!list) body = `<div class="empty-state">${MOON}<h3>${t('scanning')}</h3></div>`;
  else if (!list.length) body = `<div class="empty-state">${MOON}<h3>${t('st_empty')}</h3></div>`;
  else body = `<div class="card st-list">${[...list].sort((a, b) => (b.Enabled - a.Enabled) || (a.Desc || a.Name).localeCompare(b.Desc || b.Name)).map(x => {
    const locked = x.Admin && !admin, name = x.Desc || x.Name;
    const src = x.Source === 'hkcu' ? t('st_user') : x.Source === 'folder' || x.Source === 'common' ? t('st_folder') : t('st_all');
    return `<div class="st ${x.Enabled ? '' : 'off'}"><div class="st-ic">${esc(name.replace(/[^A-Za-z0-9]/g, '').slice(0, 2) || '?')}</div>
      <div class="tx"><b title="${esc(x.Command)}">${esc(name)}${x.Admin ? ` <span class="pill adm" title="${t('st_admin_tip')}">${t('admin')}</span>` : ''}</b><span>${esc([x.Maker, src].filter(Boolean).join(' · '))}</span></div>
      <label class="switch ${locked ? 'locked' : ''}" title="${locked ? t('st_admin_tip') : ''}"><input type="checkbox" data-st="${esc(x.Key)}" ${x.Enabled ? 'checked' : ''} ${locked ? 'disabled' : ''}><span></span></label></div>`;
  }).join('')}</div>`;
  const on = (list || []).filter(x => x.Enabled).length;
  $('#p-startup').innerHTML = `<div class="head"><div><h1>${t('st_title')}</h1><p>${t('st_sub')}</p></div>
    <div class="actions">${list ? `<span class="pill pink">${t('st_on', on, list.length)}</span>` : ''}<button class="btn icon" data-act="stRefresh" title="${t('rescan')}">${ic('refresh')}</button></div></div>
    <p class="page-note">${t('st_note')}</p>${body}${hiddenN ? `<p class="page-note" style="margin-top:14px">${t('st_hidden', hiddenN)}</p>` : ''}`;
}

async function toggleStartup(key, on) {
  try { await API().SetStartup(key, on); } catch (e) { toast(t('error', e), 'bad'); }
  await loadStartup();
}

async function scanJunk() {
  try { S.junk = await API().ScanJunk() || []; } catch { S.junk = []; }
  if (S.page === 'clean') renderClean();
}

async function cleanJunk() {
  const admin = S.sys && S.sys.Admin;
  const keys = (S.junk || []).filter(j => S.junkSel.has(j.Key) && (!j.Admin || admin)).map(j => j.Key);
  if (!keys.length) return;
  S.junkBusy = true; renderClean();
  const freed = await API().CleanJunk(keys);
  S.junkBusy = false;
  toast(t('junk_freed', fmtSize(freed)), 'ok');
  S.junk = null; renderClean(); scanJunk();
}

async function makeRestorePoint(btn) {
  btn.disabled = true;
  toast(t('restore_busy'));
  const r = await API().CreateRestorePoint();
  btn.disabled = false;
  if (r === 'ok') toast(t('restore_ok'), 'ok');
  else if (r === 'recent') toast(t('restore_recent'), 'ok');
  else if (r === 'cancelled') toast(t('restore_cancel'), 'bad');
  else toast(t('restore_err', r.replace(/^ERR\s*/, '')), 'bad');
}

// ---------- reparatie ----------
const REPAIRS = [
  { id: 'windows', icon: 'shield', confirm: 'rp_confirm_win', restart: true },
  { id: 'update', icon: 'refresh' },
  { id: 'network', icon: 'globe', confirm: 'rp_confirm_net', restart: true },
  { id: 'dns', icon: 'zap', user: true },
  { id: 'explorer', icon: 'folder', user: true },
  { id: 'icons', icon: 'palette', user: true },
  { id: 'audio', icon: 'play' },
  { id: 'bluetooth', icon: 'spark' },
  { id: 'print', icon: 'pkg' },
  { id: 'search', icon: 'search' },
  { id: 'time', icon: 'refresh' },
  { id: 'disk', icon: 'disk' },
  { id: 'store', icon: 'pkg' },
];

function renderRepair() {
  const admin = S.sys && S.sys.Admin;
  $('#p-repair').innerHTML = `<div class="head"><div><h1>${t('rp_title')}</h1><p>${t('rp_sub')}</p></div></div>
    <div class="rp-list">${REPAIRS.map(r => {
      const p = S.rpProg[r.id];
      const busy = p && p.phase === 'bezig';
      const st = p ? (p.phase === 'klaar' ? `<span class="pill ok"><i></i>${t('step_done')}</span>${r.restart ? `<span class="rp-restart">${t('rp_restart')}</span>` : ''}`
        : p.phase === 'fout' ? `<span class="pill bad"><i></i>${t('step_failed')}</span>` : '') : '';
      return `<div class="card rp ${busy ? 'busy' : ''}"><div class="rp-ic">${ic(r.icon)}</div>
        <div class="tx"><div class="rp-h"><b>${t('rp_' + r.id)}</b>${!r.user ? `<span class="pill adm">${t('admin')}</span>` : ''}<span class="rp-t">${t('rp_' + r.id + '_t')}</span></div>
          <p>${t('rp_' + r.id + '_d')}</p>
          ${busy ? `<div class="prog" style="display:flex"><div class="txt"><b>${esc(p.stepText || '')}</b><span>${p.pct >= 0 ? p.pct + '%' : ''}</span></div>
            <div class="bar ${p.pct >= 0 ? '' : 'indet'}"><i style="width:${Math.max(0, p.pct)}%"></i></div></div>` : `<div class="rp-st">${st}</div>`}</div>
        <button class="btn ${p && p.phase === 'klaar' ? '' : 'primary'} sm" data-act="repair" data-id="${r.id}" ${S.busy ? 'disabled' : ''}>${ic('wrench')}${p && p.phase === 'klaar' ? t('rp_again') : t('rp_run')}</button></div>`;
    }).join('')}</div>`;
}

async function runRepair(id) {
  const r = REPAIRS.find(x => x.id === id);
  delete S.rpProg[id];
  if (r.confirm && !(await confirmBox({ title: t('rp_' + id), text: t(r.confirm), ok: t('rp_run') }))) return;
  const msg = await API().Repair(id);
  if (msg) return toast(msg === 'busy' ? t('busy') : t('error', msg), 'bad');
  setBusy(true);
  if (!S.rpProg[id] || S.rpProg[id].phase !== 'bezig') S.rpProg[id] = { phase: 'bezig', pct: -1, stepText: t('step_starting') };
  renderRepair();
}

// ---------- privacy ----------
async function loadPrivacy() {
  try { S.pv = await API().Privacy() || []; } catch { S.pv = []; }
  if (S.page === 'privacy') renderPrivacy();
}

function pvCard(c) {
  const apps = c.Apps || [], open = S.pvOpen.has(c.ID), MAX = 6;
  const shown = open ? apps : apps.slice(0, MAX);
  return `<div class="card pv"><div class="pv-head"><div class="rp-ic">${ic(c.ID === 'webcam' ? 'desktop' : c.ID === 'microphone' ? 'spark' : 'globe')}</div>
      <div class="tx"><b>${t('pv_' + c.ID)}</b><span>${t('pv_all')}</span></div>
      <label class="switch"><input type="checkbox" data-pv="${c.ID}" data-key="" ${c.On ? 'checked' : ''}><span></span></label></div>
    ${apps.length ? `<div class="pv-apps ${c.On ? '' : 'off'}">${shown.map(x => `<div class="pv-app">
        <div class="tx"><b title="${esc(x.Path || x.Name)}">${esc(x.Name)}</b><span class="${x.InUse ? 'live' : ''}">${x.InUse ? '● ' + t('pv_inuse') : x.LastUsed ? t('pv_last', fmtDate(x.LastUsed)) : t('pv_never')}</span></div>
        <label class="switch"><input type="checkbox" data-pv="${c.ID}" data-key="${esc(x.Key)}" ${x.Allowed ? 'checked' : ''} ${c.On ? '' : 'disabled'}><span></span></label></div>`).join('')}</div>
      ${apps.length > MAX ? `<button class="pv-more" data-act="pvMore" data-id="${c.ID}">${open ? t('pv_less') : t('pv_more', apps.length)}</button>` : ''}`
      : `<div class="drv-state">${t('pv_none')}</div>`}</div>`;
}

function renderPrivacy() {
  let body;
  if (!S.pv) body = `<div class="empty-state">${MOON}<h3>${t('scanning')}</h3></div>`;
  else {
    // twee kolommen die elk doorlopen: geen gaten onder korte kaarten
    const cols = [[], []], h = [0, 0];
    for (const c of S.pv) { const i = h[0] <= h[1] ? 0 : 1; cols[i].push(c); h[i] += Math.min((c.Apps || []).length, S.pvOpen.has(c.ID) ? 99 : 6) + 3; }
    body = `<div class="pv-grid">${cols.map(col => `<div class="pv-col">${col.map(pvCard).join('')}</div>`).join('')}</div>`;
  }
  const pg = $('#p-privacy'), keep = pg.scrollTop;
  pg.innerHTML = `<div class="head"><div><h1>${t('pv_title')}</h1><p>${t('pv_sub')}</p></div>
    <div class="actions"><button class="btn icon" data-act="pvRefresh" title="${t('rescan')}">${ic('refresh')}</button></div></div>
    <p class="page-note">${t('pv_note')}</p>${body}`;
  pg.scrollTop = keep;
}

// ---------- rechtsklikmenu ----------
async function loadCtx() {
  try { S.ctx = await API().ContextMenu() || []; } catch { S.ctx = []; }
  if (S.page === 'clean') renderClean();
}

function ctxBox(bare) {
  let inner;
  if (!S.ctx) inner = `<div class="drv-state"><span class="spinner"></span>${t('scanning')}</div>`;
  else if (!S.ctx.length) inner = `<div class="drv-state ok">${ic('check')}${t('cm_none')}</div>`;
  else inner = S.ctx.map(c => `<div class="dev"><div class="tx"><b>${esc(c.Name)}</b><span class="meta">${esc(c.App || t('cm_unknown'))} · ${c.Where.split(',').map(w => t('cm_' + w)).join(', ')}</span></div>
      <label class="switch"><input type="checkbox" data-ctx="${esc(c.ID)}" ${c.Enabled ? 'checked' : ''}><span></span></label></div>`).join('');
  return `${bare ? '' : `<div class="section-title">${t('cm_title')}</div>`}<div class="card drvbox"><div class="drv-head"><div><b>${t('cm_title')}</b><span>${t('cm_sub')} ${t('cm_note')}</span></div></div>${inner}</div>`;
}

// ---------- nieuwe pc / back-up ----------
async function exportProfile() {
  try { const p = await API().ExportProfile(); if (p) toast(t('pf_saved', p), 'ok'); } catch (e) { toast(t('error', e), 'bad'); }
}

async function importProfile() {
  let p;
  try { p = await API().ImportProfile(); } catch { return toast(t('pf_bad'), 'bad'); }
  if (!p) return;
  const admin = S.sys && S.sys.Admin;
  const apps = (p.apps || []).map(id => S.apps.find(a => a.ID === id)).filter(a => a && !a.Installed);
  const tw = (p.tweaks || []).map(id => S.tweaks.find(x => x.ID === id)).filter(x => x && !x.On);
  const skip = tw.filter(x => x.Admin && !admin), doTw = tw.filter(x => !x.Admin || admin);
  if (!apps.length && !doTw.length) return toast(t('pf_nothing'), 'ok');
  const list = `<div class="un-apps">${apps.map(a => `<span class="pill pink">${esc(a.Name)}</span>`).join('')}</div>
    ${doTw.length ? `<div class="un-apps">${doTw.map(x => `<span class="pill">${esc(L(x, 'Name'))}</span>`).join('')}</div>` : ''}`;
  const ok = await confirmBox({ title: t('pf_preview'), ok: t('pf_apply'),
    text: `${t('pf_from', p.created)}<br><br><b>${t('pf_apps', apps.length)}</b> · <b>${t('pf_tweaks', doTw.length)}</b>${list}${skip.length ? `<br>${t('pf_admin_skip', skip.length)}` : ''}` });
  if (!ok) return;
  for (const x of doTw) { try { const n = await API().SetTweak(x.ID, true); S.tweaks[S.tweaks.findIndex(y => y.ID === x.ID)] = n; } catch { } }
  if (doTw.length) toast(t('pf_tweaks', doTw.length) + ' ✓', 'ok');
  if (apps.length) { goPage('apps'); install(apps.map(a => a.ID)); }
}

// ---------- Dusk zelf bijwerken ----------
async function checkSelfUpdate(manual) {
  try { S.upd = await API().CheckUpdate(); } catch { S.upd = { Err: 'offline' }; }
  if (manual && S.upd.Err) toast(t('up_offline'), 'bad');
  else if (manual && !S.upd.newer) toast(t('up_latest'), 'ok');
  if (S.upd && S.upd.newer && !manual) toast(t('up_avail', S.upd.version));
  if (S.page === 'settings') renderSettings();
  if (S.page === 'home') renderHome();
}

async function applySelfUpdate() {
  S.updPct = 0; render();
  const r = await API().ApplyUpdate();
  if (r) { S.updPct = null; render(); toast(t('up_fail', r.replace(/^ERR\s*/, '')), 'bad'); }
}

function updBanner() {
  if (!S.upd || !S.upd.newer) return '';
  const notes = LANG === 'nl' ? S.upd.notesNL : S.upd.notesEN;
  return `<div class="banner pink"><div class="ic">${ic('down')}</div><div class="tx"><b>${t('up_avail', esc(S.upd.version))}</b>${notes ? `<span>${esc(notes)}</span>` : ''}</div>
    ${S.updPct !== null ? `<span class="pill pink">${t('up_busy', S.updPct)}</span>` : `<button class="btn primary sm" data-act="selfUpdate">${ic('down')}${t('up_install')}</button>`}</div>`;
}

// ---------- opschonen ----------
function renderClean() {
  const pg = $('#p-clean');
  let body;
  if (!S.bloatLoaded) body = `<div class="empty-state">${MOON}<h3>${t('searching')}</h3>${t('searching_sub')}</div>`;
  else if (S.bloatErr) body = `<div class="banner warn"><div class="ic">${ic('alert')}</div><div class="tx"><b>${t('read_failed')}</b><span>${esc(S.bloatErr)}</span></div></div>`;
  else if (!S.bloat.length) body = `<div class="empty-state">${MOON}<h3>${t('all_clean')}</h3>${t('all_clean_sub')}</div>`;
  else body = `<div class="grid">${S.bloat.map(b => {
    const p = S.bloatProg[b.Pattern];
    const st = p ? stepText(p) : esc(L(b, 'Desc'));
    return `<div class="card bl ${S.bloatSel.has(b.Pattern) ? 'sel' : ''} ${p && p.phase === 'klaar' ? 'gone' : ''}" data-act="bsel" data-pat="${esc(b.Pattern)}">
      <div class="ico">${ic('pkg')}</div>
      <div class="tx"><div class="nm">${esc(L(b, 'Name'))}</div><div class="st">${st}</div><div class="pk">${esc(b.Pattern)}</div></div>
      <div class="check ${S.bloatSel.has(b.Pattern) ? 'on' : ''}">${ic('check')}</div></div>`;
  }).join('')}</div>`;
  const admin = S.sys && S.sys.Admin;
  const jk = S.junk;
  const jsel = (S.junk || []).filter(j => S.junkSel.has(j.Key) && (!j.Admin || admin));
  const jtotal = jsel.reduce((a, j) => a + j.Size, 0);
  const junkBox = `<div class="card drvbox junk"><div class="drv-head"><div><b>${t('junk_title')}</b><span>${t('junk_sub')}</span></div>
      <button class="btn primary sm" data-act="junkClean" ${!jk || !jtotal || S.junkBusy ? 'disabled' : ''}>${ic('trash')}${t('junk_clean', fmtSize(jtotal))}</button></div>
    ${!jk ? `<div class="drv-state"><span class="spinner"></span>${t('junk_scanning')}</div>` : jk.map(j => {
      const locked = j.Admin && !admin;
      return `<div class="drv ${S.junkSel.has(j.Key) && !locked ? 'sel' : ''} ${locked ? 'locked' : ''}" data-act="jsel" data-id="${j.Key}">
        <div class="check ${S.junkSel.has(j.Key) && !locked ? 'on' : ''}">${ic('check')}</div>
        <div class="tx"><b>${t('junk_' + j.Key)}${j.Admin ? ` <span class="pill adm">${t('junk_admin')}</span>` : ''}</b><span>${t('junk_' + j.Key + '_d')}</span></div>
        <span class="jsize">${fmtSize(j.Size)}<em>${t('files_n', j.Files)}</em></span></div>`;
    }).join('')}</div>`;
  const ctxN = S.ctx ? S.ctx.length : null, ctxOff = S.ctx ? S.ctx.filter(c => !c.Enabled).length : 0;
  const junkAll = (jk || []).filter(j => !j.Admin || admin).reduce((a, j) => a + j.Size, 0);
  const tiles = [
    ['files', 'trash', t('junk_title'), jk ? fmtSize(junkAll) : '…', t('cl_files_sub')],
    ['apps', 'pkg', t('junk_apps_title'), S.bloatLoaded ? String(S.bloat.length) : '…', t('cl_apps_sub')],
    ['menu', 'apps', t('cm_title'), ctxN === null ? '…' : String(ctxN), ctxN ? t('cl_menu_sub', ctxOff) : t('cl_menu_sub0')],
  ];
  const sec = S.cleanTab || 'files';
  const content = sec === 'files' ? junkBox
    : sec === 'apps' ? `<div class="cl-bar"><span>${t('clean_sub')}</span>${S.bloat.length ? `<button class="btn ghost sm" data-act="bAll">${t('select_all')}</button>` : ''}</div>${body}`
    : ctxBox(true);
  pg.innerHTML = `
    <div class="head"><div><h1>${t('clean_title')}</h1><p>${t('clean_page_sub')}</p></div>
      <div class="actions"><button class="btn icon" data-act="cleanRefresh" title="${t('rescan')}">${ic('refresh')}</button></div></div>
    <div class="cl-tiles">${tiles.map(([k, i, n, v, sub]) => `<button class="card cl-tile ${sec === k ? 'on' : ''}" data-cltab="${k}">
      <div class="rp-ic">${ic(i)}</div><div class="tx"><span>${n}</span><b>${v}</b><em>${sub}</em></div></button>`).join('')}</div>
    <div class="cl-body">${content}</div>`;
}

async function scanBloat() {
  try { S.bloat = await API().ScanBloat(); S.bloatErr = null; } catch (e) { S.bloatErr = String(e); }
  S.bloatLoaded = true;
  S.bloatProg = {};
  [...S.bloatSel].forEach(p => { if (!S.bloat.some(b => b.Pattern === p)) S.bloatSel.delete(p); });
  renderNav();
  render();
}

async function removeSelected() {
  const pats = [...S.bloatSel];
  const names = S.bloat.filter(b => S.bloatSel.has(b.Pattern)).map(b => esc(L(b, 'Name')));
  const ok = await confirmBox({ title: t('remove_q', pats.length), text: `${names.join(', ')}.<br><br>${t('remove_text')}`, ok: t('remove'), danger: true });
  if (!ok) return;
  const msg = await API().RemoveBloat(pats);
  if (msg) return toast(t(msg));
  setBusy(true);
  S.bloatSel.clear();
  render();
}

// ---------- instellingen ----------
const BG_PV = {
  schemering: { pv: 'background:linear-gradient(180deg,#0b0820,#2a1150 60%,#7b2c6e);', txt: '<span style="position:absolute;bottom:4px;left:4px;right:4px;color:rgba(255,140,210,.5)">.:-=+=-:..:-==-:.</span>' },
  regen: { pv: 'background:#0b0820;', txt: '1 { / 0 # | = 1 * 0 > & ; 1 [ 0' },
  golven: { pv: 'background:#0b0820;', txt: '.:=+*#*+=:. .:=+*+=:. :=+*#%#*+=:' },
  stil: { pv: 'background:linear-gradient(180deg,#0b0820,#2a1150 70%,#4a1655);', txt: '' },
};

function renderSettings() {
  const s = S.sys || {};
  $('#p-settings').innerHTML = `
    <div class="head"><div><h1>${t('settings_title')}</h1><p>${t('settings_sub')}</p></div></div>
    <div class="section-title" style="margin-top:6px">${t('language')}</div>
    <div class="card"><div class="langs">
      ${[['en', 'English'], ['nl', 'Nederlands'], ['de', 'Deutsch'], ['fr', 'Français'], ['es', 'Español']].map(([k, n]) => `<button class="langopt ${LANG === k ? 'on' : ''}" data-lang="${k}"><span class="code">${k.toUpperCase()}</span>${n}${LANG === k ? ic('check') : ''}</button>`).join('')}
    </div></div>
    <div class="section-title">${t('background')}</div>
    <div class="card"><div class="bgs">${FX.MODES.map(m => `<div class="bgopt ${S.set.bg === m ? 'on' : ''}" data-bg="${m}">
      <div class="pv" style="${BG_PV[m].pv}">${BG_PV[m].txt}</div><span>${t('bg_' + m)}</span></div>`).join('')}</div></div>
    <div class="section-title">${t('display')}</div>
    <div class="card">
      <div class="set-row"><div class="tx"><b>${t('intro_set')}</b><span>${t('intro_set_sub')}</span></div>
        <button class="btn sm" data-act="playIntro">${ic('play')}${t('play_now')}</button>
        <label class="switch"><input type="checkbox" data-set="intro" ${S.set.intro ? 'checked' : ''}><span></span></label></div>
      <div class="set-row"><div class="tx"><b>${t('calm_set')}</b><span>${t('calm_set_sub')}</span></div>
        <label class="switch"><input type="checkbox" data-set="calm" ${S.set.calm ? 'checked' : ''}><span></span></label></div>
    </div>
    <div class="section-title">${t('rights')}</div>
    <div class="card"><div class="set-row"><div class="tx"><b>${s.Admin ? t('runs_admin') : t('runs_user')}</b>
      <span>${s.Admin ? t('runs_admin_sub') : t('runs_user_sub')}</span></div>
      ${s.Admin ? `<span class="pill ok"><i></i>${t('admin')}</span>` : `<button class="btn sm" data-act="admin">${ic('shield')}${t('run_admin')}</button>`}</div></div>
    <div class="section-title">${t('pf_title')}</div>
    <div class="card"><div class="set-row"><div class="tx"><b>${t('pf_title')}</b><span>${t('pf_sub')}</span></div>
      <button class="btn sm" data-act="pfExport">${ic('down')}${t('pf_export')}</button>
      <button class="btn sm" data-act="pfImport">${ic('up')}${t('pf_import')}</button></div></div>
    <div class="section-title">${t('restore_title')}</div>
    <div class="card"><div class="set-row"><div class="tx"><b>${t('restore_btn')}</b><span>${t('restore_sub')}</span></div>
      <button class="btn sm" data-act="restore">${ic('shield')}${t('restore_btn')}</button></div></div>
    ${updBanner()}
    <div class="section-title">${t('about')}</div>
    <div class="card about"><img class="eye" src="eye.svg" alt="">${FX.logoSVG(9, 'al')}<div>
      <b style="font-size:16px">Dusk ${esc(s.Version || '')}</b>
      <p>${t('made_by')}</p>
      <div style="margin-top:10px">${S.upd && S.upd.newer ? '' : `<button class="btn ghost sm" data-act="selfCheck">${ic('refresh')}${t('up_check')}</button>`}</div>
    </div></div>`;
}

function setLang(l) {
  LANG = S.set.lang = l;
  save();
  try { API().SetLang(l); } catch { }
  applyStatic();
  renderNav();
  render();
}

// ---------- intro + achtergrond ----------
let BG;
function playIntro(first) {
  const box = $('#intro');
  box.style.display = 'block';
  box.style.opacity = '1';
  const finish = FX.intro(box.querySelector('canvas'), () => {
    box.style.display = 'none';
    if (first) showShell();
  });
  box.querySelector('.skip').onclick = finish;
  box.onclick = e => { if (e.target.tagName === 'CANVAS') finish(); };
}

function showShell() {
  const sh = $('#shell');
  sh.classList.remove('hide');
  sh.style.animation = 'none'; void sh.offsetWidth; sh.style.animation = '';
}

function applySettings() {
  BG.set(S.set.bg);
  BG.calm(S.set.calm);
  document.body.classList.toggle('calm', S.set.calm);
}

// ---------- klikken ----------
document.addEventListener('click', async e => {
  const el = e.target;
  if (el.closest('[data-score]')) return showScore();
  const nav = el.closest('[data-page]');
  if (nav) return goPage(nav.dataset.page);
  const gt = el.closest('[data-goto]');
  if (gt) return goPage(gt.dataset.goto);
  const ps = el.closest('[data-psort]');
  if (ps) { S.progSort = ps.dataset.psort; return renderPrograms(); }
  const vw = el.closest('[data-view]');
  if (vw) { S.view = vw.dataset.view; $('#p-apps').scrollTop = 0; return renderAppGrid(); }
  const cl = el.closest('[data-cltab]');
  if (cl) { S.cleanTab = cl.dataset.cltab; return renderClean(); }
  const jump = el.closest('[data-jump]');
  if (jump) {
    const sec = $('#sec-' + jump.dataset.jump), pg = $('#p-apps');
    if (sec) pg.scrollTo({ top: sec.offsetTop - $('#catBar').offsetHeight - 14, behavior: S.set.calm ? 'auto' : 'smooth' });
    return;
  }
  const bg = el.closest('[data-bg]');
  if (bg) { S.set.bg = bg.dataset.bg; save(); applySettings(); return renderSettings(); }
  const lg = el.closest('[data-lang]');
  if (lg) return setLang(lg.dataset.lang);
  if (el.closest('[data-goto-oem]')) { const b = $('#oemBox'); if (b) b.scrollIntoView({ behavior: 'smooth', block: 'start' }); return; }
  const url = el.closest('[data-url]');
  if (url) return window.runtime.BrowserOpenURL(url.dataset.url);
  const info = el.closest('[data-info]');
  if (info) { const id = info.dataset.info; S.twOpen.has(id) ? S.twOpen.delete(id) : S.twOpen.add(id); return renderTweaks(); }
  const twRow = el.closest('.tw-row');
  if (twRow && !el.closest('label.switch')) {
    const inp = twRow.querySelector('input[data-tweak]');
    if (inp && !inp.disabled) { inp.checked = !inp.checked; return toggleTweak(inp.dataset.tweak, inp.checked); }
    return;
  }
  const a = el.closest('[data-act]');
  if (!a) {
    const card = el.closest('.app');
    if (card && !card.classList.contains('busy') && !el.closest('button')) toggleSel(card.dataset.id);
    return;
  }
  const id = a.dataset.id;
  switch (a.dataset.act) {
    case 'sel': return toggleSel(id);
    case 'one': return install([id]);
    case 'updAll': {
      const ids = S.apps.filter(x => x.Available).map(x => x.ID);
      if (S.page !== 'apps') goPage('apps');
      return install(ids);
    }
    case 'installSel': return install([...S.sel].filter(i => { const x = S.apps.find(y => y.ID === i); return x && (!x.Installed || x.Available); }));
    case 'clearSel': S.sel.clear(); return render();
    case 'bsel': {
      const p = a.dataset.pat;
      if (S.bloatProg[p]) return;
      S.bloatSel.has(p) ? S.bloatSel.delete(p) : S.bloatSel.add(p);
      renderClean(); return renderSel();
    }
    case 'bAll': S.bloat.forEach(b => S.bloatSel.add(b.Pattern)); renderClean(); return renderSel();
    case 'clearBloat': S.bloatSel.clear(); renderClean(); return renderSel();
    case 'removeSel': return removeSelected();
    case 'rescanBloat': S.bloatLoaded = false; renderClean(); return scanBloat();
    case 'hwRefresh': S.hw = null; S.hwErr = null; S.drvStat = null; renderHardware(); await loadHardware(); return checkDrivers();
    case 'drvScan': case 'drvCheck': return checkDrivers();
    case 'osel': {
      if (S.oemProg[id]) return;
      S.oemSel.has(id) ? S.oemSel.delete(id) : S.oemSel.add(id);
      renderHardware(); return renderSel();
    }
    case 'jsel': {
      const j = (S.junk || []).find(x => x.Key === id);
      if (!j || (j.Admin && !(S.sys && S.sys.Admin))) return;
      S.junkSel.has(id) ? S.junkSel.delete(id) : S.junkSel.add(id);
      return renderClean();
    }
    case 'junkClean': return cleanJunk();
    case 'cleanRefresh': S.junk = null; S.ctx = null; S.bloatLoaded = false; renderClean(); scanJunk(); loadCtx(); return scanBloat();
    case 'wgInstall': {
      const msg = await API().InstallWinget();
      if (msg === 'ok') { S.sys.Winget = true; return scanApps(); }
      if (msg) return toast(t(msg));
      setBusy(true); S.wg = { phase: 'bezig', step: 'w_deps' }; return renderAppGrid();
    }
    case 'repair': return runRepair(id);
    case 'pvMore': S.pvOpen.has(id) ? S.pvOpen.delete(id) : S.pvOpen.add(id); return renderPrivacy();
    case 'pvRefresh': S.pv = null; renderPrivacy(); return loadPrivacy();
    case 'pfExport': return exportProfile();
    case 'pfImport': return importProfile();
    case 'selfUpdate': return applySelfUpdate();
    case 'selfCheck': return checkSelfUpdate(true);
    case 'psel': { if (S.progProg[id]) return; S.progSel.has(id) ? S.progSel.delete(id) : S.progSel.add(id); renderPrograms(); return renderSel(); }
    case 'clearProg': S.progSel.clear(); renderPrograms(); return renderSel();
    case 'progRemove': return uninstallFlow('progs', [...S.progSel]);
    case 'appRemove': return uninstallFlow('apps', [...S.sel].filter(i => { const x = S.apps.find(y => y.ID === i); return x && x.Installed && x.Cat !== 'runtime'; }));
    case 'gpuInstall': {
      const msg = await API().InstallGPUDriver(a.dataset.name, a.dataset.dl);
      if (msg) return toast(msg === 'busy' ? t('busy') : t('error', msg), 'bad');
      setBusy(true); S.gpuProg[a.dataset.name] = { phase: 'bezig', pct: 0, step: 'downloading' };
      return renderHardware();
    }
    case 'wuOne': {
      const msg = await API().InstallDrivers([id]);
      if (msg) return toast(t(msg));
      setBusy(true); S.drvProg[id] = { phase: 'bezig', pct: -1, step: 'installing' };
      return renderHardware();
    }
    case 'amdOpen':
      if (await API().OpenVendorApp('amd')) return toast(t('amd_opening'));
      return window.runtime.BrowserOpenURL(VENDOR.amd.url);
    case 'intelCheck': {
      const dsa = S.apps.find(x => x.ID === 'Intel.IntelDriverAndSupportAssistant');
      if (dsa && dsa.Installed) { await API().OpenVendorApp('intel'); toast(t('intel_opening')); return window.runtime.BrowserOpenURL(VENDOR.intel.url); }
      toast(t('intel_installing'));
      S.afterInstall = 'intel';
      return install(['Intel.IntelDriverAndSupportAssistant']);
    }
    case 'oemRun': {
      const msg = await API().RunOEMTool(S.hw.OEM.Brand);
      if (msg === 'notool') return install([S.hw.OEM.Tool]);
      if (msg) return toast(t(msg));
      setBusy(true); return renderHardware();
    }
    case 'stRefresh': S.st = null; renderStartup(); return loadStartup();
    case 'restore': return makeRestorePoint(a);
    case 'dsel': {
      if (S.drvProg[id]) return;
      S.drvSel.has(id) ? S.drvSel.delete(id) : S.drvSel.add(id);
      renderHardware(); return renderSel();
    }
    case 'clearDrv': S.drvSel.clear(); S.oemSel.clear(); renderHardware(); return renderSel();
    case 'drvInstall': {
      // eerst de pakketten van de maker (HP), dan Windows Update
      const oem = (S.drvStat && S.drvStat.OEM || []).filter(u => S.oemSel.has(u.ID));
      const wu = [...S.drvSel];
      if (oem.length) {
        const msg = await API().InstallOEM(oem);
        if (msg) return toast(t(msg));
        S.pendingWU = wu;
      } else {
        const msg = await API().InstallDrivers(wu);
        if (msg) return toast(t(msg));
      }
      setBusy(true);
      S.drvSel.clear(); S.oemSel.clear();
      return render();
    }
    case 'devmgr': return API().OpenDeviceManager().catch(e => toast(t('error', e), 'bad'));
    case 'ignore': S.set.ignored = [...new Set([...(S.set.ignored || []), a.dataset.name])]; save(); renderNav(); return renderHardware();
    case 'showIgnored': S.set.ignored = []; save(); renderNav(); return renderHardware();
    case 'explorer':
      await API().RestartExplorer();
      S.explorerPending = false;
      toast(t('explorer_done'), 'ok');
      return renderTweaks();
    case 'admin':
      try { await API().RunAsAdmin(); } catch { toast(t('admin_failed'), 'bad'); }
      return;
    case 'playIntro': return playIntro(false);
  }
});

function toggleSel(id) {
  if (!id) return;
  S.sel.has(id) ? S.sel.delete(id) : S.sel.add(id);
  updateCard(id);
  document.querySelectorAll('#appSections .grid').forEach(g => g.classList.toggle('selecting', S.sel.size > 0));
  renderSel();
}

document.addEventListener('change', e => {
  const el = e.target;
  if (el.dataset.tweak) return toggleTweak(el.dataset.tweak, el.checked);
  if (el.dataset.st) return toggleStartup(el.dataset.st, el.checked);
  if (el.dataset.pv) {
    return API().SetPrivacy(el.dataset.pv, el.dataset.key, el.checked).then(loadPrivacy).catch(e => { toast(t('error', e), 'bad'); loadPrivacy(); });
  }
  if (el.dataset.ctx) {
    return API().SetContextItem(el.dataset.ctx, el.checked).then(loadCtx).catch(e => { toast(t('error', e), 'bad'); loadCtx(); });
  }
  if (el.dataset.set) { S.set[el.dataset.set] = el.checked; save(); applySettings(); }
});

document.addEventListener('keydown', e => {
  if (e.key === 'Escape') {
    if ($('#drawer').classList.contains('on')) $('#drawer').classList.remove('on');
    else if (S.sel.size) { S.sel.clear(); render(); }
  }
  if (e.ctrlKey && e.key === 'f' && S.page === 'apps') { e.preventDefault(); $('#q').focus(); }
});

// venster
$('#btnMin').onclick = () => window.runtime.WindowMinimise();
$('#btnMax').onclick = () => window.runtime.WindowToggleMaximise();
$('#titlebar').ondblclick = e => { if (!e.target.closest('button')) window.runtime.WindowToggleMaximise(); };
$('#btnClose').onclick = async () => {
  if (S.busy && !(await confirmBox({ title: t('closing_busy'), text: t('closing_busy_text'), ok: t('close_anyway'), danger: true }))) return;
  window.runtime.Quit();
};
$('#btnLog').onclick = () => $('#drawer').classList.toggle('on');
$('#logClose').onclick = () => $('#drawer').classList.remove('on');
$('#logClear').onclick = () => { $('#log').innerHTML = ''; };

// Terug in Dusk na een paar minuten (bv. drivers bijgewerkt in een ander programma)?
// Dan de huidige pagina stil opnieuw controleren.
window.addEventListener('focus', async () => {
  if (S.busy || Date.now() - (S.lastRefresh || 0) < 120000) return;
  S.lastRefresh = Date.now();
  if (S.page === 'hardware' && S.hw && !S.drvChecking) { await loadHardware(); checkDrivers(); }
  else if (S.page === 'apps' && !S.scanning) scanApps();
  else if (S.page === 'startup') loadStartup();
  else if (S.page === 'clean') { scanJunk(); scanBloat(); }
});

// ---------- gebeurtenissen van Go ----------
function hookEvents() {
  on('log', log);
  on('app', p => {
    S.prog[p.id] = p;
    if (p.phase === 'klaar' || p.phase === 'fout') {
      const a = S.apps.find(x => x.ID === p.id);
      if (a && p.phase === 'klaar') {
        if (p.step === 'removed') { a.Installed = false; a.Version = ''; a.Available = ''; }
        else { a.Installed = true; if (a.Available) a.Version = a.Available; a.Available = ''; }
        renderNav();
      }
      if (a) toast(p.phase === 'klaar' ? t('app_done', a.Name) : t('app_failed', a.Name, p.detail), p.phase === 'klaar' ? 'ok' : 'bad');
      if (p.phase === 'fout' && a && a.Installed) S.failed.add(p.id);
      setTimeout(() => { delete S.prog[p.id]; updateCard(p.id); }, 50);
      return updateCard(p.id, p.phase === 'klaar' ? 'done' : 'fail');
    }
    updateCard(p.id);
  });
  on('apps', list => { S.apps = list; S.appsLoaded = true; renderNav(); render(); });
  on('bloat', p => { S.bloatProg[p.id] = p; if (S.page === 'clean') renderClean(); });
  on('driver', p => { S.drvProg[p.id] = p; if (S.page === 'hardware') renderHardware(); });
  on('oem', p => { S.oemProg[p.id] = p; if (S.page === 'hardware') renderHardware(); });
  on('winget', p => {
    S.wg = p;
    if (p.phase === 'klaar') { setBusy(false); S.wg = null; S.sys.Winget = true; toast(t('wg_done'), 'ok'); scanApps(); }
    if (p.phase === 'fout') { setBusy(false); toast(t('wg_fail'), 'bad'); setTimeout(() => { S.wg = null; render(); }, 4000); }
    if (S.page === 'apps') renderAppGrid();
  });
  on('repair', p => {
    const cur = S.rpProg[p.id] || { pct: -1, stepText: '' };
    cur.phase = p.phase;
    if (p.step && p.step.startsWith('r_')) { cur.stepText = t('step_' + p.step); cur.pct = -1; }
    else if (!p.step && p.pct >= 0) cur.pct = p.pct;
    if (p.step === 'failed') cur.stepText = t('step_failed');
    S.rpProg[p.id] = cur;
    if (S.page === 'repair') renderRepair();
  });
  on('selfupdate', pct => { S.updPct = pct; if (S.page === 'settings' || S.page === 'home') render(); });
  on('prog', p => { S.progProg[p.id] = p; if (S.page === 'apps' && S.view === 'pc') renderPrograms(); });
  on('gpu', p => { S.gpuProg[p.id] = p; if (S.page === 'hardware') renderHardware(); });
  on('oemtool', p => { S.oemToolProg[p.id] = p; });
  on('done', d => {
    setBusy(false);
    if (d.kind === 'bloat') {
      toast(d.fail ? t('removed_fail', d.ok, d.fail) : t('removed_n', d.ok), d.fail ? 'bad' : 'ok');
      setTimeout(scanBloat, 1200);
    } else if (d.kind === 'repair') {
      toast(d.fail ? t('step_failed') : t('step_done'), d.fail ? 'bad' : 'ok');
      if (S.page === 'repair') renderRepair();
      return;
    } else if (d.kind === 'uninstall') {
      S.progProg = {}; S.progs = null;
      showUninstallResult(d);
      if (S.page === 'apps') renderAppGrid();
      return;
    } else if (d.kind === 'drivers' || d.kind === 'oem') {
      toast(d.fail ? t('drivers_fail', d.ok, d.fail) : t('drivers_done', d.ok), d.fail ? 'bad' : 'ok');
      if (d.kind === 'oem' && S.pendingWU && S.pendingWU.length) {
        const wu = S.pendingWU; S.pendingWU = null;
        API().InstallDrivers(wu).then(m => { if (!m) setBusy(true); });
        return render();
      }
      S.gpuProg = {};
      setTimeout(async () => { await loadHardware(); checkDrivers(); }, 1500);
    } else {
      toast(d.fail ? t('some_failed', d.ok, d.fail) : t('all_done'), d.fail ? 'bad' : 'ok');
      if (S.afterInstall === 'intel' && !d.fail) { S.afterInstall = null; API().OpenVendorApp('intel'); window.runtime.BrowserOpenURL(VENDOR.intel.url); toast(t('intel_opening')); }
      const failed = [...S.failed];
      S.failed.clear();
      if (failed.length) { render(); return install(failed, true); }
    }
    render();
  });
}

// ---------- start ----------
async function start() {
  applyStatic();
  $('#logo').innerHTML = FX.logoSVG(8, 'sl');
  BG = FX.background($('#bg'));
  applySettings();
  hookEvents();
  try { API().SetLang(LANG); } catch { }
  renderNav();
  let noIntro = false;
  try { noIntro = await API().NoIntro(); } catch { }
  if (S.set.intro && !noIntro) playIntro(true);
  else { $('#intro').style.display = 'none'; showShell(); }
  goPage('home');

  try {
    S.sys = await API().System();
    $('#tbVer').textContent = 'v' + S.sys.Version;
    const user = S.sys.User || '?';
    const parts = user.split(/[\s._-]+/).filter(Boolean);
    $('#whoAv').textContent = (parts.length > 1 ? parts[0][0] + parts[1][0] : user.slice(0, 2)).toUpperCase();
    $('#whoAdm').innerHTML = S.sys.Admin ? `<span class="admdot" title="${t('admin')}">${ic('shield')}</span>` : '';
    $('#whoName').textContent = user;
    $('#whoPc').textContent = S.sys.PC || '';
    $('.who').title = `${user} · ${S.sys.PC}`;
    fitText($('#whoName'), 13.5, 10);
    fitText($('#whoPc'), 12, 9);
  } catch (e) { toast(t('sys_failed', e), 'bad'); }
  render();
  API().Tweaks().then(x => { S.tweaks = x; S.tweaksLoaded = true; renderNav(); render(); });
  scanBloat();
  loadHardware().then(() => checkDrivers()); // op de achtergrond, ook voor de pc-score
  try { if (await API().JustUpdated()) toast(t('up_done', S.sys ? S.sys.Version : ''), 'ok'); } catch { }
  setTimeout(() => checkSelfUpdate(false), 4000);
  if (S.sys && S.sys.Winget) scanApps();
  else { S.appsLoaded = true; try { S.apps = await API().ScanApps(); } catch { } render(); }
}

const ERRS = [];
window.addEventListener('error', e => ERRS.push(String(e.message)));
window.addEventListener('unhandledrejection', e => ERRS.push('promise: ' + String(e.reason)));

async function selfTest() {
  try { if (!(await API().SelfTest())) return; } catch { return; }
  API().InUse(['OBSProject.OBSStudio', 'Discord.Discord']).then(r => { window.INUSE = r; }).catch(e => { window.INUSE = String(e); });
  setTimeout(() => {
    const r = { errors: ERRS, inuse: window.INUSE, page: S.page, lang: LANG, apps: S.apps.length, appsLoaded: S.appsLoaded, tweaks: S.tweaks.length, bloat: S.bloat.length,
      hw: S.hw, sys: S.sys, navItems: document.querySelectorAll('.nav').length, w: innerWidth, h: innerHeight };
    API().SelfTestDone(JSON.stringify(r, null, 1));
  }, 15000);
}

window.addEventListener('DOMContentLoaded', () => {
  // wacht tot Wails klaar is (window.go bestaat)
  const t0 = Date.now();
  (function wait() {
    if (window.go && window.go.main && window.runtime) { selfTest(); return start(); }
    if (Date.now() - t0 > 5000) return document.body.insertAdjacentHTML('beforeend', `<p style="position:fixed;top:50%;width:100%;text-align:center">${t('start_failed')}</p>`);
    setTimeout(wait, 30);
  })();
});
