// Maakt het oog-icoon in alle maten (groot = met details, klein = simpel en dik).
const { chromium } = require('playwright');
const fs = require('fs');
function svg(simple) {
  // iris-strepen
  let rays = '';
  if (!simple) for (let i = 0; i < 48; i++) {
    const a = i / 48 * Math.PI * 2, r1 = 22, r2 = 41 - (i % 3) * 4;
    rays += `<line x1="${128 + Math.cos(a) * r1}" y1="${132 + Math.sin(a) * r1}" x2="${128 + Math.cos(a) * r2}" y2="${132 + Math.sin(a) * r2}" stroke="rgba(255,220,255,${0.18 + (i % 4) * 0.06})" stroke-width="2"/>`;
  }
  const eye = 'M24 132 C70 62 186 62 232 132 C186 202 70 202 24 132 Z';
  return `<svg width="256" height="256" viewBox="0 0 256 256" xmlns="http://www.w3.org/2000/svg"><defs>
  <linearGradient id="bg" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#120a2e"/><stop offset=".6" stop-color="#341262"/><stop offset="1" stop-color="#7a2470"/></linearGradient>
  <radialGradient id="gl" cx=".5" cy="1" r=".75"><stop offset="0" stop-color="#ff6fb5" stop-opacity=".5"/><stop offset="1" stop-color="#ff6fb5" stop-opacity="0"/></radialGradient>
  <linearGradient id="wh" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#f6efff"/><stop offset="1" stop-color="#d6c6f5"/></linearGradient>
  <radialGradient id="ir" cx=".5" cy=".5" r=".5"><stop offset=".3" stop-color="#ff5fb4"/><stop offset=".62" stop-color="#c04cf0"/><stop offset=".92" stop-color="#7c3aed"/><stop offset="1" stop-color="#3b0f7a"/></radialGradient>
  <radialGradient id="pu" cx=".5" cy=".5" r=".5"><stop offset="0" stop-color="#ff7ac8"/><stop offset=".35" stop-color="#5a1050"/><stop offset=".8" stop-color="#0a0418"/></radialGradient>
  <clipPath id="ec"><path d="${eye}"/></clipPath>
  <filter id="b"><feGaussianBlur stdDeviation="${simple ? 4 : 8}"/></filter></defs>
  <rect x="6" y="6" width="244" height="244" rx="58" fill="url(#bg)"/><rect x="6" y="6" width="244" height="244" rx="58" fill="url(#gl)"/>
  ${simple ? '' : '<rect x="7.5" y="7.5" width="241" height="241" rx="56.5" fill="none" stroke="rgba(255,255,255,.13)" stroke-width="3"/>'}
  <path d="${eye}" fill="#c084fc" filter="url(#b)" opacity=".75"/>
  <path d="${eye}" fill="url(#wh)"/>
  <g clip-path="url(#ec)">
    <circle cx="128" cy="132" r="${simple ? 50 : 46}" fill="url(#ir)"/>${rays}
    <circle cx="128" cy="132" r="${simple ? 22 : 20}" fill="url(#pu)"/>
    ${simple ? '' : '<circle cx="128" cy="132" r="2.6" fill="#fff"/><path d="M128 124v16M120 132h16" stroke="#fff" stroke-width="1.6" opacity=".8"/>'}
    <circle cx="${simple ? 146 : 144}" cy="${simple ? 114 : 116}" r="${simple ? 10 : 8}" fill="#fff" opacity=".9"/>
    <path d="M24 132 C70 62 186 62 232 132" fill="none" stroke="#2a0f55" stroke-width="${simple ? 10 : 7}" opacity=".35"/>
  </g>
  <path d="${eye}" fill="none" stroke="#e9d5ff" stroke-width="${simple ? 9 : 5}"/>
  ${simple ? '' : '<g fill="#fff"><circle cx="196" cy="46" r="4"/><circle cx="216" cy="70" r="2.5" opacity=".75"/><circle cx="52" cy="210" r="2.5" opacity=".6"/></g>'}
  </svg>`;
}
(async () => {
  const b = await chromium.launch();
  for (const [size, simple] of [[256, 0], [64, 0], [48, 1], [32, 1], [24, 1], [16, 1]]) {
    const p = await b.newPage({ viewport: { width: size, height: size } });
    await p.setContent(`<body style="margin:0;background:transparent"><img src="data:image/svg+xml;base64,${Buffer.from(svg(simple)).toString('base64')}" width="${size}" height="${size}"></body>`);
    await p.screenshot({ path: `winres/eye${size}.png`, omitBackground: true });
    await p.close();
  }
  fs.writeFileSync('frontend/eye.svg', svg(0));
  fs.writeFileSync('frontend/eye-small.svg', svg(1));
  await b.close();
})();
