package script

func GetInjectionScript() string {
	return `
(function(){
  // --- 0. Zoom Manager (hotkey only: Ctrl +/-/0, Ctrl+Wheel) ---
  var zoomLevels = [0.7,0.8,0.85,0.9,0.95,1.0,1.05,1.1,1.2,1.35,1.5];
  var zoomIndex = 5;
  try{
    var saved = localStorage.getItem('mdeck_zoom_factor');
    if(saved){ var f=parseFloat(saved); for(var i=0;i<zoomLevels.length;i++){ if(Math.abs(zoomLevels[i]-f)<0.01){ zoomIndex=i; break; } } }
  }catch(e){}
  function showToast(msg){
    var t=document.getElementById('mdeck-toast');
    if(t) t.remove();
    var el=document.createElement('div');
    el.id='mdeck-toast';
    el.textContent=msg;
    el.style.cssText='position:fixed;bottom:22px;left:50%;transform:translateX(-50%);background:rgba(20,22,28,.92);color:#e8e8e5;padding:8px 14px;border-radius:8px;border:1px solid rgba(255,255,255,.08);font-size:12px;z-index:99999;backdrop-filter:blur(10px);pointer-events:none;font-family:Inter,system-ui,sans-serif;';
    document.body.appendChild(el);
    setTimeout(function(){ el.remove(); }, 1600);
  }
  function applyZoom(silent){
    var f=zoomLevels[zoomIndex];
    if('zoom' in document.body.style){ document.documentElement.style.zoom = f; }
    else { document.documentElement.style.transform='scale('+f+')'; document.documentElement.style.transformOrigin='top left'; document.documentElement.style.width=(100/f)+'%'; document.documentElement.style.height=(100/f)+'%'; }
    // keep titlebar physical 38px (counter-scale)
    var bar=document.getElementById('mdeck-titlebar');
    if(bar){
      if('zoom' in bar.style) bar.style.zoom = 1/f;
      else { bar.style.transform='scale('+(1/f)+')'; bar.style.transformOrigin='top left'; bar.style.width=(100*f)+'%'; }
    }
    try{ localStorage.setItem('mdeck_zoom_factor', String(f)); }catch(e){}
    if(!silent) showToast('Zoom '+Math.round(f*100)+'%');
    setTimeout(function(){ window.dispatchEvent(new Event('resize')); }, 80);
  }
  window.__mdeck_triggerZoomIn = function(){ if(zoomIndex<zoomLevels.length-1){ zoomIndex++; applyZoom(); } };
  window.__mdeck_triggerZoomOut = function(){ if(zoomIndex>0){ zoomIndex--; applyZoom(); } };
  window.__mdeck_triggerResetZoom = function(){ zoomIndex=5; applyZoom(); };
  setTimeout(function(){ applyZoom(true); }, 120);

  // --- 1. Titlebar + Resize Handles ---
  if(document.getElementById('mdeck-titlebar')){ return; }
  var styleId='mdeck-titlebar-style';
  if(!document.getElementById(styleId)){
    var s=document.createElement('style');
    s.id=styleId;
    s.textContent=[
      "#mdeck-titlebar{position:fixed;top:0;left:0;right:0;height:38px;display:flex;align-items:center;justify-content:space-between;padding:0 12px;background:rgba(12,14,18,.88);backdrop-filter:blur(16px) saturate(1.1);color:#e6e6e6;border-bottom:1px solid rgba(255,255,255,.06);z-index:2147483640;--wails-draggable:drag;user-select:none;font-family:Inter,system-ui,sans-serif;transition:background .18s,border-color .18s;}",
      ".mdeck-traffic{display:flex;gap:8px;--wails-draggable:no-drag}",
      ".mdeck-dot{width:12px;height:12px;border-radius:50%;border:none;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:transform .12s;}",
      ".mdeck-dot:active{transform:scale(.9)}",
      ".mdeck-dot.close{background:#ff5f56;box-shadow:inset 0 0 0 1px #e0443e} .mdeck-dot.min{background:#ffbd2e;box-shadow:inset 0 0 0 1px #dea123} .mdeck-dot.max{background:#27c93f;box-shadow:inset 0 0 0 1px #1aab29}",
      ".mdeck-traffic .mdeck-dot svg{opacity:0;transition:opacity .15s} .mdeck-traffic:hover .mdeck-dot svg{opacity:.85}",
      ".mdeck-center{position:absolute;left:50%;transform:translateX(-50%);font-size:13px;font-weight:650;letter-spacing:-.2px;opacity:.92;pointer-events:none;display:flex;align-items:center;gap:8px}",
      ".mdeck-right{display:flex;gap:4px;--wails-draggable:no-drag;align-items:center}",
      ".mdeck-btn{width:30px;height:26px;border-radius:7px;border:none;background:transparent;color:#8b9099;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:background .12s,color .12s;}",
      ".mdeck-btn:hover{background:rgba(255,255,255,.08);color:#fff}",
      ".mdeck-btn:active{background:rgba(255,255,255,.12)}",
      ".mdeck-resize{position:fixed;z-index:2147483647;background:transparent;--wails-draggable:no-drag}",
      ".mdeck-resize.top{top:0;left:10px;right:10px;height:5px;cursor:n-resize} .mdeck-resize.bottom{bottom:0;left:10px;right:10px;height:5px;cursor:s-resize}",
      ".mdeck-resize.left{top:10px;bottom:10px;left:0;width:5px;cursor:w-resize} .mdeck-resize.right{top:10px;bottom:10px;right:0;width:5px;cursor:e-resize}",
      ".mdeck-resize.tl{top:0;left:0;width:10px;height:10px;cursor:nw-resize} .mdeck-resize.tr{top:0;right:0;width:10px;height:10px;cursor:ne-resize}",
      ".mdeck-resize.bl{bottom:0;left:0;width:10px;height:10px;cursor:sw-resize} .mdeck-resize.br{bottom:0;right:0;width:10px;height:10px;cursor:se-resize}",
      "body.mdeck-has-titlebar{padding-top:38px !important;transition:padding-top .2s;}",
    ].join("\n");
    (document.head||document.documentElement).appendChild(s);
  }
  function invoke(msg){
    try{ if(window.chrome&&window.chrome.webview&&window.chrome.webview.postMessage){ window.chrome.webview.postMessage(msg); return; } }catch(e){}
    try{ if(window.webkit&&window.webkit.messageHandlers&&window.webkit.messageHandlers.external){ window.webkit.messageHandlers.external.postMessage(msg); return; } }catch(e){}
    try{ if(window.WailsInvoke){ window.WailsInvoke(msg); return; } }catch(e){}
  }

  // Traffic-light buttons use Wails runtime bridge directly (works on WebKitGTK & WebView2)
  function doClose(){ invoke('Q'); try{ window.go.main.App.Close(); }catch(e){ try{ window.runtime.Quit(); }catch(e2){} } }
  function doMin(){ invoke('Wm'); try{ window.go.main.App.Minimise(); }catch(e){ try{ window.runtime.WindowMinimise(); }catch(e2){} } }
  function doMax(){ invoke('Wt'); try{ window.go.main.App.ToggleMaximise(); }catch(e){ try{ window.runtime.WindowToggleMaximise(); }catch(e2){} } }

  var bar=document.createElement('div'); bar.id='mdeck-titlebar';
  bar.innerHTML='<div class="mdeck-traffic">'
    +'<button class="mdeck-dot close" title="Close"><svg width=6 height=6 viewBox="0 0 6 6"><line x1=1 y1=1 x2=5 y2=5 stroke="#4d0000" stroke-width=1.2 stroke-linecap=round/><line x1=5 y1=1 x2=1 y2=5 stroke="#4d0000" stroke-width=1.2 stroke-linecap=round/></svg></button>'
    +'<button class="mdeck-dot min" title="Minimize"><svg width=6 height=6 viewBox="0 0 6 6"><line x1=1 y1=3 x2=5 y2=3 stroke="#664400" stroke-width=1.2 stroke-linecap=round/></svg></button>'
    +'<button class="mdeck-dot max" title="Maximize"><svg width=8 height=8 viewBox="0 0 12 12"><rect x=1.5 y=1.5 width=9 height=9 fill=none stroke="#0a4a0a" stroke-width=1.2/></svg></button>'
    +'</div><div class="mdeck-center">MDeck</div><div class="mdeck-right">'
    +'<button class="mdeck-btn" title="New Window (Ctrl N)" data-act="newwin"><svg width=13 height=13 viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M3 7v10a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-4"/><path d="M8 3h8l3 3v6a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2z"/></svg></button>'
    +'</div>';
  // wire traffic lights via addEventListener (onclick inline can break under CSP / webview quirks)
  bar.querySelector('.mdeck-dot.close').addEventListener('click', function(e){ e.stopPropagation(); doClose(); });
  bar.querySelector('.mdeck-dot.min').addEventListener('click', function(e){ e.stopPropagation(); doMin(); });
  bar.querySelector('.mdeck-dot.max').addEventListener('click', function(e){ e.stopPropagation(); doMax(); });
  var nwBtn = bar.querySelector('[data-act=newwin]');
  if(nwBtn) nwBtn.addEventListener('click', function(e){ e.stopPropagation(); try{ window.go.main.App.NewWindow(); }catch(err){} });

  var last=0;
  bar.addEventListener('mousedown', function(e){
    if(e.target.closest('button')) return;
    if(e.button!==0) return;
    var now=Date.now(); if(now-last<280){ e.preventDefault(); doMax(); return; } last=now;
    invoke('drag');
  });
  bar.addEventListener('dblclick', function(e){ if(e.target.closest('button')) return; doMax(); });
  document.documentElement.appendChild(bar);
  document.body.classList.add('mdeck-has-titlebar');

  [['top'],['bottom'],['left'],['right'],['tl'],['tr'],['bl'],['br']].forEach(function(k){
    var cls=k[0]; if(document.querySelector('.mdeck-resize.'+cls)) return;
    var h=document.createElement('div'); h.className='mdeck-resize '+cls;
    h.addEventListener('mousedown', function(e){ if(e.button===0){ e.preventDefault(); var map={top:'n-resize',bottom:'s-resize',left:'w-resize',right:'e-resize',tl:'nw-resize',tr:'ne-resize',bl:'sw-resize',br:'se-resize'}; invoke('resize:'+map[cls]); }});
    document.documentElement.appendChild(h);
  });

  // --- 2. Global shortcuts (zoom via hotkey only) ---
  if(!window.__mdeck_shortcuts){
    window.__mdeck_shortcuts=true;
    window.addEventListener('keydown', function(e){
      var mod=e.ctrlKey||e.metaKey;
      if(mod && (e.key==='='||e.key==='+'||e.code==='Equal'||e.code==='NumpadAdd')){ e.preventDefault(); window.__mdeck_triggerZoomIn(); return; }
      if(mod && (e.key==='-'||e.key==='_'||e.code==='Minus'||e.code==='NumpadSubtract')){ e.preventDefault(); window.__mdeck_triggerZoomOut(); return; }
      if(mod && (e.key==='0'||e.code==='Digit0'||e.code==='Numpad0')){ e.preventDefault(); window.__mdeck_triggerResetZoom(); return; }
      if(mod && e.key.toLowerCase()==='n'){ e.preventDefault(); try{ window.go.main.App.NewWindow(); }catch(err){} return; }
    }, true);
    window.addEventListener('wheel', function(e){
      if(e.ctrlKey){ e.preventDefault(); if(e.deltaY<0) window.__mdeck_triggerZoomIn(); else window.__mdeck_triggerZoomOut(); }
    }, {passive:false});
  }
})();
`
}
