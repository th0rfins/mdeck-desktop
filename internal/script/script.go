package script

func GetInjectionScript() string {
	return `
(function(){
  // inject frameless traffic-light titlebar + 8 resize handles, like Gotion
  // but simplified for MDeck: dark #090a0d, traffic lights left, actions right
  if (document.getElementById('mdeck-titlebar')) { syncTheme(); return; }

  var styleId='mdeck-titlebar-style';
  if(!document.getElementById(styleId)){
    var s=document.createElement('style');
    s.id=styleId;
    s.textContent=[
      "#mdeck-titlebar{position:fixed;top:0;left:0;right:0;height:38px;display:flex;align-items:center;justify-content:space-between;padding:0 10px;background:#0c0e12;color:#e6e6e6;border-bottom:1px solid rgba(255,255,255,.06);z-index:2147483640;--wails-draggable:drag;user-select:none;font-family:Inter,system-ui,sans-serif;}",
      ".mdeck-traffic{display:flex;gap:8px;--wails-draggable:no-drag}",
      ".mdeck-dot{width:12px;height:12px;border-radius:50%;border:1px solid rgba(0,0,0,.2);cursor:pointer;display:flex;align-items:center;justify-content:center}",
      ".mdeck-dot.close{background:#ff5f56;border-color:#e0443e} .mdeck-dot.min{background:#ffbd2e;border-color:#dea123} .mdeck-dot.max{background:#27c93f;border-color:#1aab29}",
      ".mdeck-traffic .mdeck-dot svg{opacity:0} .mdeck-traffic:hover .mdeck-dot svg{opacity:.9}",
      ".mdeck-center{position:absolute;left:50%;transform:translateX(-50%);font-size:13px;font-weight:600;opacity:.9;pointer-events:none}",
      ".mdeck-right{display:flex;gap:4px;--wails-draggable:no-drag}",
      ".mdeck-btn{width:28px;height:28px;border-radius:6px;border:none;background:transparent;color:#8b9099;cursor:pointer;display:flex;align-items:center;justify-content:center}",
      ".mdeck-btn:hover{background:rgba(255,255,255,.08);color:#fff}",
      ".mdeck-resize{position:fixed;z-index:2147483647;background:transparent;--wails-draggable:no-drag}",
      ".mdeck-resize.top{top:0;left:10px;right:10px;height:6px;cursor:n-resize} .mdeck-resize.bottom{bottom:0;left:10px;right:10px;height:6px;cursor:s-resize}",
      ".mdeck-resize.left{top:10px;bottom:10px;left:0;width:6px;cursor:w-resize} .mdeck-resize.right{top:10px;bottom:10px;right:0;width:6px;cursor:e-resize}",
      ".mdeck-resize.tl{top:0;left:0;width:10px;height:10px;cursor:nw-resize} .mdeck-resize.tr{top:0;right:0;width:10px;height:10px;cursor:ne-resize}",
      ".mdeck-resize.bl{bottom:0;left:0;width:10px;height:10px;cursor:sw-resize} .mdeck-resize.br{bottom:0;right:0;width:10px;height:10px;cursor:se-resize}",
      // push MDeck web content below titlebar when inside launcher (wails://)
      "body.mdeck-has-titlebar{padding-top:38px !important}",
      // browser-tab polish for LabDeck TabBar when loaded under mdeck wrapper
      ".mdeck-tabbar-polish{padding-top:38px}",
    ].join("\n");
    (document.head||document.documentElement).appendChild(s);
  }

  function invoke(msg){
    try{ if(window.chrome&&window.chrome.webview&&window.chrome.webview.postMessage){ window.chrome.webview.postMessage(msg); return; } }catch(e){}
    try{ if(window.webkit&&window.webkit.messageHandlers&&window.webkit.messageHandlers.external){ window.webkit.messageHandlers.external.postMessage(msg); return; } }catch(e){}
  }

  var bar=document.createElement('div'); bar.id='mdeck-titlebar';
  bar.innerHTML='<div class="mdeck-traffic">'
    +'<button class="mdeck-dot close" title="Close" onclick="window.go&&window.go.main&&window.go.main.App&&window.go.main.App.Close()"><svg width=6 height=6 viewBox="0 0 6 6"><line x1=1 y1=1 x2=5 y2=5 stroke="#4d0000" stroke-width=1.2 stroke-linecap=round/><line x1=5 y1=1 x2=1 y2=5 stroke="#4d0000" stroke-width=1.2 stroke-linecap=round/></svg></button>'
    +'<button class="mdeck-dot min" title="Minimize" onclick="window.go&&window.go.main&&window.go.main.App&&window.go.main.App.Minimise()"><svg width=6 height=6 viewBox="0 0 6 6"><line x1=1 y1=3 x2=5 y2=3 stroke="#664400" stroke-width=1.2 stroke-linecap=round/></svg></button>'
    +'<button class="mdeck-dot max" title="Maximize" onclick="window.go&&window.go.main&&window.go.main.App&&window.go.main.App.ToggleMaximise()"><svg width=8 height=8 viewBox="0 0 12 12"><rect x=1.5 y=1.5 width=9 height=9 fill=none stroke="#0a4a0a" stroke-width=1.2/></svg></button>'
    +'</div><div class="mdeck-center">MDeck</div><div class="mdeck-right">'
    +'<button class="mdeck-btn" title="Back" onclick="history.back()">‹</button>'
    +'<button class="mdeck-btn" title="Reload" onclick="location.reload()">↻</button>'
    +'<button class="mdeck-btn" title="New Window" onclick="window.go&&window.go.main&&window.go.main.App&&window.go.main.App.NewWindow()">＋</button>'
    +'</div>';
  // drag & dblclick toggle maximize
  var last=0;
  bar.addEventListener('mousedown', function(e){
    if(e.target.closest('button')) return;
    if(e.button!==0) return;
    var now=Date.now(); if(now-last<300){ e.preventDefault(); if(window.go&&window.go.main&&window.go.main.App) window.go.main.App.ToggleMaximise(); return; } last=now;
    invoke('drag');
  });
  bar.addEventListener('dblclick', function(e){ if(e.target.closest('button')) return; if(window.go&&window.go.main&&window.go.main.App) window.go.main.App.ToggleMaximise(); });
  document.documentElement.appendChild(bar);
  document.body.classList.add('mdeck-has-titlebar');

  // 8 resize handles
  [['top'],['bottom'],['left'],['right'],['tl'],['tr'],['bl'],['br']].forEach(function(k){
    var cls=k[0]; if(document.querySelector('.mdeck-resize.'+cls)) return;
    var h=document.createElement('div'); h.className='mdeck-resize '+cls;
    h.addEventListener('mousedown', function(e){ if(e.button===0){ e.preventDefault(); var map={top:'n-resize',bottom:'s-resize',left:'w-resize',right:'e-resize',tl:'nw-resize',tr:'ne-resize',bl:'sw-resize',br:'se-resize'}; invoke('resize:'+map[cls]); }});
    document.documentElement.appendChild(h);
  });

  function syncTheme(){
    var isDark = !document.body.classList.contains('light');
    var b=document.getElementById('mdeck-titlebar');
    if(b) b.style.background = isDark ? '#0c0e12' : '#f6f5f4';
  }
  syncTheme();
  setInterval(syncTheme, 800);
})();
`
}
