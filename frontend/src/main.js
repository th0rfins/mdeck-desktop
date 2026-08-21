// MDeck Hub — bookmark manager before startup
// Uses Wails bindings: window.go.main.App.{GetBookmarks,SaveBookmark,DeleteBookmark,Navigate,NewWindow,OpenExternal}

const colors = ["violet","teal","blue","grape","orange","cyan"];

function el(html){ const t=document.createElement('template'); t.innerHTML=html.trim(); return t.content.firstElementChild; }

async function getBookmarks(){
  try{ if(window.go?.main?.App?.GetBookmarks) return await window.go.main.App.GetBookmarks(); }catch(e){}
  // fallback for browser preview
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

function render(list){
  const app=document.getElementById('app');
  app.innerHTML='';
  const header=el(`<div class="header"><div><div class="title">MDeck Hub</div><div class="hint">Pilih mesin untuk konek. Bookmark tersimpan lokal — bisa warna & buka multi-window seperti terminal.</div></div><button class="btn" id="newWinHint">New Window: Ctrl+N</button></div>`);
  app.appendChild(header);

  const form=el(`<div class="form">
    <input id="label" placeholder="Label (HK1 Production)" />
    <input id="url" placeholder="https://hk1.projectpop.xyz" />
    <select id="color">${colors.map(c=>`<option value="${c}">${c}</option>`).join('')}</select>
    <button class="btn primary" id="addBtn">Save bookmark</button>
  </div>`);
  app.appendChild(form);

  const grid=el(`<div class="grid"></div>`);
  list.forEach(b=>{
    const card=el(`<div class="card" style="border-left:3px solid var(--c)">
      <div style="display:flex;justify-content:space-between;align-items:center"><span class="badge" style="background:var(--c);color:white">${b.color}</span><span class="hint">${b.id.slice(0,8)}</span></div>
      <div class="label" style="margin-top:10px">${b.label}</div>
      <div class="url">${b.url}</div>
      <div class="row">
        <button class="btn primary connect">Connect</button>
        <button class="btn newwin">New Window</button>
        <button class="btn ghost del">Delete</button>
      </div>
    </div>`);
    card.style.setProperty('--c', b.color==='violet'?'#7048e8': b.color==='teal'?'#0ca678': b.color==='blue'?'#3578e5':'#868e96');
    card.querySelector('.connect').onclick=()=>connect(b.url);
    card.querySelector('.newwin').onclick=()=>newWindow(b.url);
    card.querySelector('.del').onclick=async()=>{ await deleteBookmark(b.id); render(await getBookmarks()); };
    card.ondblclick=()=>connect(b.url);
    grid.appendChild(card);
  });
  app.appendChild(grid);

  if(list.length===0){
    app.appendChild(el(`<div class="hint" style="margin-top:18px;text-align:center">Belum ada bookmark — tambahkan hk1/hk2 di form atas.</div>`));
  }

  form.querySelector('#addBtn').onclick=async()=>{
    const label=form.querySelector('#label').value.trim();
    const url=form.querySelector('#url').value.trim();
    const color=form.querySelector('#color').value;
    if(!label||!url) return alert('Label & URL wajib');
    const b={id: Math.random().toString(36).slice(2,10), label, url, color};
    const next=await saveBookmark(b);
    form.querySelector('#label').value=''; form.querySelector('#url').value='';
    render(next);
  };
}

getBookmarks().then(render);
document.addEventListener('keydown', e=>{
  if((e.ctrlKey||e.metaKey)&& e.key.toLowerCase()==='n'){ e.preventDefault(); newWindow(''); }
});
