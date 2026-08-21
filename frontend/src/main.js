// MDeck Hub — polished bookmark manager
const colors = {
  violet:{bg:'#7048e8', glow:'rgba(112,72,232,.55)'},
  teal:{bg:'#0ca678', glow:'rgba(12,166,120,.45)'},
  blue:{bg:'#3578e5', glow:'rgba(53,120,229,.5)'},
  grape:{bg:'#ae3ec9', glow:'rgba(174,62,201,.45)'},
  orange:{bg:'#ff6b35', glow:'rgba(255,107,53,.45)'},
  cyan:{bg:'#15aabf', glow:'rgba(21,170,191,.45)'},
};
const colorKeys = Object.keys(colors);

function el(html){ const t=document.createElement('template'); t.innerHTML=html.trim(); return t.content.firstElementChild; }

async function getBookmarks(){
  try{ if(window.go?.main?.App?.GetBookmarks) return await window.go.main.App.GetBookmarks(); }catch(e){}
  const raw=localStorage.getItem('mdeck_bookmarks');
  if(raw) try{ return JSON.parse(raw); }catch(e){}
  return [
    {id:'hk1', label:'HK1 Production', url:'https://hk1.projectpop.xyz', color:'violet'},
    {id:'hk2', label:'HK2 Staging', url:'https://hk2.projectpop.xyz', color:'teal'},
  ];
}
async function saveBookmark(b){
  if(window.go?.main?.App?.SaveBookmark) return await window.go.main.App.SaveBookmark(b);
  const list=await getBookmarks();
  const idx=list.findIndex(x=>x.id===b.id);
  if(idx>=0) list[idx]=b; else list.push(b);
  localStorage.setItem('mdeck_bookmarks', JSON.stringify(list));
  return list;
}
async function deleteBookmark(id){
  if(window.go?.main?.App?.DeleteBookmark) return await window.go.main.App.DeleteBookmark(id);
  const list=(await getBookmarks()).filter(x=>x.id!==id);
  localStorage.setItem('mdeck_bookmarks', JSON.stringify(list));
  return list;
}
function connect(url){
  if(window.go?.main?.App?.Navigate) window.go.main.App.Navigate(url);
  else window.location.href=url;
}
function newWindow(url){
  if(window.go?.main?.App?.NewWindow) window.go.main.App.NewWindow(url);
  else window.open(url, '_blank');
}

function colorFor(c){ return colors[c] || colors.violet; }

function render(list, filter=''){
  const app=document.getElementById('app');
  app.innerHTML='';
  const filtered = filter ? list.filter(b=> (b.label+b.url).toLowerCase().includes(filter.toLowerCase())) : list;

  // hero
  const hero = el(`
    <div class="hero">
      <div class="hero-left">
        <div class="hero-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="3"/><path d="M8 9l3 3-3 3"/><path d="M13 15h4"/></svg>
        </div>
        <div>
          <div style="font-weight:800;letter-spacing:-.6px;font-size:16px">MDeck Hub</div>
          <div class="hero-sub">Pilih mesin untuk konek — bookmark dengan warna, buka multi-window seperti terminal. <span class="kbd">Ctrl N</span> new window · <span class="kbd">Ctrl ±</span> zoom</div>
        </div>
      </div>
      <div class="hero-actions">
        <div class="search">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="11" cy="11" r="7"/><path d="M20 20l-3.5-3.5"/></svg>
          <input id="search" placeholder="Search hk1, projectpop..." />
        </div>
      </div>
    </div>
  `);
  app.appendChild(hero);

  // form
  const form = el(`
    <div class="form">
      <div class="form-field">
        <label>Label</label>
        <input id="label" placeholder="HK1 Production" />
      </div>
      <div class="form-field" style="flex:1.6">
        <label>URL</label>
        <input id="url" placeholder="https://hk1.projectpop.xyz" />
      </div>
      <div class="form-field" style="max-width:140px">
        <label>Warna</label>
        <select id="color">${colorKeys.map(c=>`<option value="${c}">${c}</option>`).join('')}</select>
      </div>
      <button class="btn primary" id="addBtn" style="height:38px;align-self:flex-end">Save bookmark</button>
    </div>
  `);
  app.appendChild(form);

  // grid
  const grid = el(`<div class="grid"></div>`);
  if(filtered.length===0){
    grid.appendChild(el(`<div class="empty" style="grid-column:1/-1">
      <h3>Belum ada mesin</h3>
      <div class="hint">Tambahkan bookmark di form atas — contoh: hk1.projectpop.xyz</div>
    </div>`));
  } else {
    filtered.forEach((b, idx)=>{
      const col = colorFor(b.color);
      const card = el(`
        <div class="card" style="--c:${col.bg};--c-glow:${col.glow};animation:fadeIn .45s ease ${idx*40}ms both">
          <div class="card-top">
            <span class="badge" style="background:rgba(255,255,255,.06);border:1px solid rgba(255,255,255,.08);color:var(--text)"><i style="background:${col.bg};box-shadow:0 0 10px ${col.glow}"></i>${b.color}</span>
            <span class="card-id">${b.id.slice(0,8)}</span>
          </div>
          <div class="card-label" title="${b.label}">${b.label}</div>
          <div class="card-url" title="${b.url}">${b.url}</div>
          <div class="card-row">
            <button class="btn primary connect" style="flex:1">Connect</button>
            <button class="btn newwin">New Window</button>
            <button class="btn ghost del" title="Delete">✕</button>
          </div>
          <div class="card-dot" style="background:${col.bg}"></div>
        </div>
      `);
      card.querySelector('.connect').onclick=()=>connect(b.url);
      card.querySelector('.newwin').onclick=()=>newWindow(b.url);
      card.querySelector('.del').onclick=async(e)=>{ e.stopPropagation(); if(confirm('Hapus '+b.label+'?')){ const next=await deleteBookmark(b.id); render(next, hero.querySelector('#search').value); } };
      card.ondblclick=()=>connect(b.url);
      // click card = connect
      card.addEventListener('click', (e)=>{ if(e.target.closest('button')) return; connect(b.url); });
      grid.appendChild(card);
    });
  }
  app.appendChild(grid);

  // footer hint
  app.appendChild(el(`<div class="hint" style="margin-top:18px;display:flex;gap:10px;flex-wrap:wrap">
    <span>Tip: <span class="kbd">Double click</span> card untuk connect · <span class="kbd">Shift+drag</span> di terminal untuk tmux selection</span>
    <span style="margin-left:auto;opacity:.7">${list.length} bookmark · ${filtered.length} shown</span>
  </div>`));

  // events
  const search = hero.querySelector('#search');
  search.value = filter;
  search.addEventListener('input', ()=> render(list, search.value));
  form.querySelector('#addBtn').onclick=async()=>{
    const label=form.querySelector('#label').value.trim();
    const url=form.querySelector('#url').value.trim();
    const color=form.querySelector('#color').value;
    if(!label||!url) return;
    if(!/^https?:\/\//.test(url)) return alert('URL harus https://...');
    const b={id: Math.random().toString(36).slice(2,9), label, url, color};
    const next=await saveBookmark(b);
    form.querySelector('#label').value=''; form.querySelector('#url').value='';
    render(next, search.value);
  };
  // enter to save
  form.addEventListener('keydown', e=>{ if(e.key==='Enter'){ e.preventDefault(); form.querySelector('#addBtn').click(); } });
}

getBookmarks().then(list=>render(list));
document.addEventListener('keydown', e=>{
  if((e.ctrlKey||e.metaKey)&& e.key.toLowerCase()==='n'){ e.preventDefault(); newWindow(''); }
});
