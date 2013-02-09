/*!

   Flowplayer v5.0.0 (Thursday, 27. September 2012 02:43PM) | flowplayer.org/license

*/
!function() { 

// auto-install (any video tag with parent .flowplayer)
$(function() {
   if (typeof $.fn.flowplayer == 'function') {
      $("video").parent(".flowplayer").flowplayer();
   }
});

var instances = [],
   extensions = [],
   iPhone = /iPhone/i.test(navigator.userAgent);


/* flowplayer()  */
window.flowplayer = function(fn) {
   return iPhone ? 0 :
      $.isFunction(fn) ? extensions.push(fn) :
      typeof fn == 'number' || fn === undefined ? instances[fn || 0] :
      $(fn).data("flowplayer");
};

$.extend(flowplayer, {

   version: '5.0.0',

   engine: {},

   conf: {

      debug: false,

      // true = forced playback
      disabled: false,

      // first engine to try
      engine: 'html5',

      // keyboard shortcuts
      keyboard: true,

      // default aspect ratio
      ratio: 9 / 16,

      rtmp: 0,

      splash: false,

      swf: "http://releases.flowplayer.org/5.0.0/flowplayer.swf",

      // initial volume level
      volume: 0.8,

      // http://www.whatwg.org/specs/web-apps/current-work/multipage/the-video-element.html#error-codes
      errors: [

         // video exceptions
         '',
         'Video loading aborted',
         'Network error',
         'Video not properly encoded',
         'Video file not found',

         // player exceptions
         'Unsupported video',
         'Skin not found'
      ]

   }

});

// iPhone completely lacks JavaScript API
if (iPhone) return;

// shortcut
$.fn.tc = $.fn.toggleClass;

// jQuery plugin
$.fn.flowplayer = function(opts, callback) {

   if (typeof opts == 'string') opts = { swf: opts }
   if ($.isFunction(opts)) { callback = opts; opts = {} }

   return this.each(function() {

      // already installed --> skip
      if ($(this).data("flowplayer")) return;

      // private variables
      var root = $(this),
         conf = $.extend({}, flowplayer.conf, opts, root.data()),
         videoTag = $("video", root),
         initialTypes = [],
         savedVolume,
         engine;


      /*** API ***/
      var api = {

         // properties
         conf: conf,
         volumeLevel: conf.volume,
         video: null,

         // states
         splash: false,
         ready: false,
         paused: false,
         playing: false,
         loading: false,
         muted: false,
         disabled: false,
         finished: false,

         // methods
         load: function(video, callback) {

            if (api.error || api.loading || api.disabled) return;

            root.trigger("load", [api, video, engine]);

            // callback
            if ($.isFunction(video)) callback = video;
            if (callback) root.one("ready", callback);

            return api;
         },

         pause: function(fn) {
            if (api.ready && !api.seeking && !api.disabled && !api.loading) {
               engine.pause();
               api.one("pause", fn);
            }
            return api;
         },

         resume: function() {

            if (api.ready && api.paused && !api.disabled) {
               engine.resume();

               // Firefox (+others?) does not fire "resume" after finish
               if (api.finished) {
                  root.trigger("resume");
                  api.finished = false;
               }
            }

            return api;
         },

         toggle: function() {
            return api.ready ? api.paused ? api.resume() : api.pause() : api.load();
         },

         seek: function(time, callback) {
            if (api.ready) {
               engine.seek(Math.min(Math.max(time, 0), api.video.duration));
               if (callback) root.one("seek", callback);
            }
            return api;
         },

         volume: function(level) {
            if (api.ready && level != api.volumeLevel) engine.volume(Math.min(Math.max(level, 0), 1));
            return api;
         },

         stop: function() {
            if (api.ready) engine.stop();
            return api;
         },

         unload: function() {
            if (!root.hasClass("is-embedding")) {
               if (conf.splash) {
                  engine.unload();
                  api.trigger("unload");
               } else {
                  api.stop();
               }
            }
            return api;
         }

      };

      /* togglers */
      $.each(['disable', 'mute'], function(i, key) {
         api[key] = function() {
            return api.trigger(key);
         };
      });

      /* event binding / unbinding */
      $.each(['bind', 'one', 'unbind'], function(i, key) {
         api[key] = function(type, fn) {
            root[key](type, fn);
            return api;
         };
      });

      api.trigger = function(event, arg) {
         root.trigger(event, [api, arg]);
         return api;
      };


      /*** Behaviour ***/

      root.bind("boot", function() {

         // conf
         $.each(['autoplay', 'loop', 'preload', 'poster'], function(i, key) {
            var val = videoTag.attr(key);
            if (val !== undefined) conf[key] = val ? val : true;
         });

         // splash
         if (conf.splash || root.hasClass("is-splash")) {
            api.splash = conf.splash = conf.autoplay = true;
            root.addClass("is-splash");
         }

         if (conf.poster) delete conf.autoplay;

         // extensions
         $.each(extensions, function(i) {
            this(api, root);
         });

         // 1. use the configured engine
         engine = flowplayer.engine[conf.engine];
         if (engine) engine = engine(api, root);

         if (engine) {
            api.engine = conf.engine;

         // 2. failed -> try another
         } else {
            delete flowplayer.engine[conf.engine];

            $.each(flowplayer.engine, function(name, impl) {
               engine = this(api, root);
               if (engine) api.engine = name;
               return false;
            });
         }

         // no engine
         if (!engine) return root.trigger("error", [api, { code: 5 }]);

         // start
         conf.splash ? api.unload() : api.load();

         // disabled
         if (conf.disabled) api.disable();

         // initial callback
         root.one("ready", callback);

         // instances
         instances.push(api);


      }).bind("load", function(e, api, video) {

         // unload others
         if (conf.splash) {
            $(".flowplayer").filter(".is-ready, .is-loading").not(root).each(function() {
               var api = $(this).data("flowplayer");
               if (api.conf.splash) api.unload();
            });
         }

         // loading
         api.loading = true;


      }).bind("ready unload", function(e) {
         var ready = e.type == "ready";
         root.tc("is-splash", !ready).tc("is-ready", ready);
         api.ready = ready;
         api.splash = !ready;

         function noLoad() {
            root.removeClass("is-loading")
            api.loading = false;
         }

         // load
         if (ready) {
            api.volume(conf.volume);

            if (conf.autoplay) {
               root.one("resume", noLoad);

            } else {
               api.trigger("pause");
               noLoad();
            }

         // unload
         } else {
            api.video.time = 0;
            if (conf.splash) videoTag.remove();
         }

      }).bind("mute", function(e) {
         var flag = api.muted = !api.muted;
         if (flag) savedVolume = api.volumeLevel;
         api.volume(flag ? 0 : savedVolume);


      }).bind("volume", function(e, api, level) {
         api.volumeLevel = Math.round(level * 100) / 100;

         if (api.muted && api.volumeLevel) {
            root.removeClass("is-muted");
            api.muted = false;
         }


      }).bind("beforeseek seek", function(e) {
         api.seeking = e.type == "beforeseek";
         root.tc("is-seeking", api.seeking);


      }).bind("ready pause resume unload finish", function(e) {

         // PAUSED: pause / finish / first-frame / poster
         api.paused = /pause|finish|unload/.test(e.type) || e.type == "ready" && !conf.autoplay;

         // the opposite
         api.playing = !api.paused;

         // CSS classes
         root.tc("is-paused", api.paused).tc("is-playing", api.playing);

         // sanity check
         if (!api.load.ed) api.pause();

      }).bind("disable", function(){
         api.disabled = !api.disabled;

      }).bind("finish", function(e) {
         api.finished = true;

      }).bind("error", function() {
         videoTag.remove();
      });

      // boot
      root.trigger("boot", [api, root]).data("flowplayer", api);

   });

};

/* The most minimal Flash embedding */

try {

   var ie = $.browser.msie,
      ver = ie ? new ActiveXObject("ShockwaveFlash.ShockwaveFlash").GetVariable('$version') :
         navigator.plugins["Shockwave Flash"].description;

   ver = ver.split(/\D+/);
   if (!ie) ver = ver.slice(1);

   $.support.flashVideo = ver[0] > 9 || ver[0] == 9 && ver[2] >= 115;

} catch (err) {

}


// movie required in opts
function embed(swf, flashvars) {

   var id = "_" + Math.random(),
      tag = '<object class="fp-engine" id="' + id+ '" name="' + id + '" ';

   tag += $.browser.msie ? 'classid="clsid:D27CDB6E-AE6D-11cf-96B8-444553540000">' :
      ' data="' + swf + '" type="application/x-shockwave-flash">';

   var opts = {
      width: "100%",
      height: "100%",
      allowscriptaccess: "always",
      wmode: "opaque",
      quality: "high",
      flashvars: "",
      movie: swf,
      name: id
   };

   // flashvars
   $.each(flashvars, function(key, value) {
      opts.flashvars += key + "=" + value + "&";
   });

   // parameters
   $.each(opts, function(key, value) {
      tag += '<param name="' + key + '" value="'+ value +'"/>';
   });

   tag += "</object>";

   return $(tag);
}


// Flash is buggy allover
if (window.attachEvent) {
   window.attachEvent("onbeforeunload", function() {
      __flash_savedUnloadHandler = __flash_unloadHandler = function() {};
   });
}


flowplayer.engine.flash = function(player, root) {

   var conf = player.conf,
      video = player.video,
      callbackId,
      objectTag,
      api;

   function pick(sources) {
      for (var i = 0, source; i < sources.length; i++) {
         source = sources[i];
         if (/mp4|flv/.test(source.type)) return source;
      }
   }

   // not supported
   if (!$.support.flashVideo || !pick(video.sources)) return;

   // ok
   $("video", root).remove();

   var engine = {

      load: function(video) {

         var source = pick(video.sources);

         // absolute url
         var url = $("<a/>").attr("href", source.src)[0];
         source.url = url.href;

         if (conf.rtmp) {
            url = url.pathname;
            if (url[0] == "/") url = url.slice(1);
            if (source.type == "mp4") url = "mp4:" + url;
            source.url = url;
         }

         if (api) {
            api.__play(source.url);

         } else {

            callbackId = "fp" + ("" + Math.random()).slice(3, 15);

            var opts = {
               hostname: conf.embedded ? conf.hostname : top.location.hostname,
               url: source.url,
               callback: "$."+ callbackId
            };

            // optional conf
            $.each(['key', 'autoplay', 'preload', 'poster', 'rtmp', 'loop'], function(i, key) {
               if (conf[key]) opts[key] = conf[key];
            });

            objectTag = embed(conf.swf, opts);

            objectTag.prependTo(root);

            api = objectTag[0];

            // listen
            $[callbackId] = function(type, arg) {

               if (conf.debug && type != "status") console.log("--", type, arg);

               var event = $.Event(type),
                  video = player.video;

               switch (type) {

                  // RTMP sends a lot of finish events in vain
                  case "finish": if (conf.rtmp) return;
                  case "ready": arg = $.extend(video, arg); break;
                  case "click": event.flash = true; break;
                  case "keydown": event.which = arg; break;
                  case "buffered": video.buffered = true; break;
                  case "seek": video.time = arg; break;

                  case "status":
                     if (!video.time || arg.time > video.time) {
                        video.time = arg.time;
                        player.trigger("progress", arg);
                     }

                     if (arg.buffer < video.bytes) {
                        video.buffer = arg.buffer / video.bytes * video.duration;
                        player.trigger("buffer", arg);
                     }
                     break;
               }

               // add some delay to that player is truly ready after an event
               setTimeout(function() { player.trigger(event, arg); }, 1)

            };

         }

         return source;

      },

      unload: function() {
         $("object", root).remove();
         api = 0;
      }

   };

   $.each("pause,resume,seek,volume,stop,unload".split(","), function(i, name) {

      engine[name] = function(arg) {
         if (api) {

            if (name == 'seek') {

               // started
               if (player.video.time) {
                  player.trigger("beforeseek");

               // not started (TODO: simplify)
               } else {
                  engine.resume();
                  player.seeking = false;
                  player.trigger("resume");
                  return;
               }

            }

            if (name == 'unload') {
               api.__unload();
               delete $[callbackId];
               objectTag.remove();
               api = 0;

            } else if (arg === undefined) {
               api["__" + name]();

            } else {
               api["__" + name](arg);
            }

         }
      };

   });

   return engine;

};


var VIDEO = $('<video/>')[0];

   // HTML5 --> Flowplayer event
var EVENTS = {

   // fired
   ended: 'finish',
   pause: 'pause',
   play: 'resume',
   progress: 'buffer',
   timeupdate: 'progress',
   volumechange: 'volume',
   seeking: 'beforeseek',
   seeked: 'seek',
   // abort: 'resume',

   // not fired
   loadeddata: 'ready',
   // loadedmetadata: 0,
   // canplay: 0,

   // error events
   // load: 0,
   // emptied: 0,
   // empty: 0,
   error: 'error',
   dataunavailable: 'error'

};

$.support.video = !!VIDEO.canPlayType;


flowplayer.engine.html5 = function(player, root) {

   var videoTag = $("video", root),
      conf = player.conf,
      timer,
      api;


   function pick(video) {
      for (var i = 0, source; i < video.sources.length; i++) {
         source = video.sources[i];
         if (VIDEO.canPlayType("video/" + source.type).replace("no", '')) return source;
      };
   }

   // not supported
   if (!$.support.video || !pick(player.video)) return;

   // ok
   videoTag.addClass("fp-engine").removeAttr("controls");

   return {

      load: function(video) {

         var source = pick(video);

         if (conf.splash && !api) {
            videoTag = $("<video/>", {
               src: source.src,
               type: 'video/' + source.type,
               autoplay: 'autoplay',
               'class': 'fp-engine'
            }).prependTo(root);

            if (conf.loop) videoTag.attr("loop", "loop");

            api = videoTag[0];

         } else {
            api = videoTag[0];

            // change of clip
            if (video.src && api.src != video.src) api.src = source.src;
            api.load();
         }

         listen(api);

         return source;
      },

      pause: function() {
         api.pause();
      },

      resume: function() {
         api.play();
      },

      seek: function(time) {
         try {
            api.currentTime = time;
         } catch (ignored) {}
      },

      // seek(0) && pause() && display poster
      stop: function() {
         api.currentTime = 0;
         setTimeout(function() { api.load(); }, 100);
      },

      volume: function(level) {
         api.volume = level;
      },

      unload: function() {
         $("video", root).remove();
         api = 0;
      }

   };

   function listen(api) {

      // listen only once
      if (api.listening) return; api.listening = true;

      $.each(EVENTS, function(type, flow) {

        api.addEventListener(type, function(e) {

            if (conf.debug && !/progress/.test(flow)) console.log(type, "->", flow, e);

            // no events if player not ready
            if (!player.ready && !/ready|error/.test(flow) || !flow) { return; }

            var event = $.Event(flow), video = player.video, arg;

            switch (flow) {

               case "ready":

                  arg = $.extend(video, {
                     duration: api.duration,
                     width: api.videoWidth,
                     height: api.videoHeight,
                     url: api.currentSrc
                  });

                  try {
                     video.seekable = api.seekable && api.seekable.end(null);

                  } catch (ignored) {}

                  // buffer
                  timer = timer || setInterval(function() {

                     try {
                        video.buffer = api.buffered.end(null);

                     } catch (ignored) {}

                     if (video.buffer) {
                        if (video.buffer < video.duration) {
                           player.trigger("buffer", e);

                        } else if (!video.buffered) {
                           video.buffered = true;
                           player.trigger("buffer", e).trigger("buffered", e);
                           clearInterval(timer);
                           timer = 0;
                        }
                     }

                  }, 250);

                  break;

               case "progress": case "seek":
                  // Safari can give negative times. add rounding
                  arg = video.time = Math.round(Math.max(api.currentTime, 0) * 100) / 100;;
                  break;

               case "volume":
                  arg = api.volume;
                  break;

               case "error":
                  arg = (e.srcElement || e.originalTarget).error;
            }

            player.trigger(event, arg);

         }, false);

      });

   }

};
var TYPE_RE = /.(\w{3,4})$/i;

function parseSource(el) {
   var type = el.attr("type"), src = el.attr("src");
   return { src: src, type: type ? type.replace("video/", "") : src.split(TYPE_RE)[1] };
}

/* Resolves video object from initial configuration and from load() method */
flowplayer(function(api, root) {

   var videoTag = $("video", root),
      initialSources = [];

   // initial video
   $("source", videoTag).each(function() {
      initialSources.push(parseSource($(this)));
   });

   if (!initialSources.length) initialSources.push(parseSource(videoTag));

   api.video = { sources: initialSources };

   // a new video is loaded
   api.bind("load", function(e, api, video, engine) {

      video = video || api.video;

      if ($.isArray(video)) {

         video = { sources: $.map(video, function(el) {
            var type = Object.keys(el)[0];
            el.type = type;
            el.src = el[type];
            delete el[type];
            return el;
         })};

      } else if (typeof video == 'string') {

         video = { src: video, sources: [] };

         $.each(initialSources, function(i, source) {
            video.sources.push({
               type: source.type,
               src: video.src.replace(TYPE_RE, "") + "." + source.type
            });
         });

      }

      api.video = $.extend(video, engine.load(video));

   });

});
/* A minimal jQuery Slider plugin with all goodies */

// skip IE policies
document.ondragstart = function () { return false; };


// execute function every <delay> ms
$.throttle = function(fn, delay) {
   var locked;

   return function () {
      if (!locked) {
         fn.apply(this, arguments);
         locked = 1;
         setTimeout(function () { locked = 0; }, delay);
      }
   };
};


$.fn.slider = function() {

   return this.each(function() {

      var root = $(this),
         doc = $(document),
         progress = root.children(":last"),
         disabled,
         offset,
         width,
         height,
         vertical,
         size,
         maxValue,
         max,

         /* private */
         calc = function() {
            offset = root.offset();
            width = root.width();
            height = root.height();

            /* exit from fullscreen can mess this up.*/
            // vertical = height > width;

            size = vertical ? height : width;
            max = toDelta(maxValue);
         },

         fire = function(value) {
            if (!disabled && value != api.value && (!maxValue || value < maxValue)) {
               root.trigger("slide", [ value ]);
               api.value = value;
            }
         },

         mousemove = function(e) {
            var delta = vertical ? e.pageY - offset.top : e.pageX - offset.left;
            delta = Math.max(0, Math.min(max || size, delta));

            var value = delta / size;
            if (vertical) value = 1 - value;
            return move(value, 0, true);
         },

         move = function(value, speed) {
            if (speed === undefined) { speed = 0; }
            var to = (value * 100) + "%";

            if (!maxValue || value <= maxValue)
               progress.stop().animate(vertical ? { height: to } : { width: to }, speed, "linear");

            return value;
         },

         toDelta = function(value) {
            return Math.max(0, Math.min(size, vertical ? (1 - value) * height : value * width));
         },

         /* public */
         api = {

            max: function(value) {
               maxValue = value;
            },

            disable: function(flag) {
               disabled = flag;
            },

            slide: function(value, speed, fireEvent) {
               calc();
               if (fireEvent) fire(value);

               move(value, speed);
            }

         };

      calc();

      // bound dragging into document
      root.data("api", api).bind("mousedown.sld", function(e) {

         e.preventDefault();

         if (!disabled) {

            // begin --> recalculate. allows dynamic resizing of the slider
            var delayedFire = $.throttle(fire, 100);
            calc();
            api.dragging = true;
            fire(mousemove(e));

            doc.bind("mousemove.sld", function(e) {
               e.preventDefault();
               delayedFire(mousemove(e));

            }).one("mouseup", function() {
               api.dragging = false;
               doc.unbind("mousemove.sld");
            });

         }

      });

   });

};

function zeropad(val) {
   val = parseInt(val, 10);
   return val >= 10 ? val : "0" + val;
}

// display seconds in hh:mm:ss format
function format(sec) {

   sec = sec || 0;

   var h = Math.floor(sec / 3600),
       min = Math.floor(sec / 60);

   sec = sec - (min * 60);

   if (h >= 1) {
      min -= h * 60;
      return h + "h:" + zeropad(min); // + ":" + zeropad(sec);
   }

   return zeropad(min) + ":" + zeropad(sec);
}


// detect animation support
$.support.animation = (function() {
   var vendors = ['','Webkit','Moz','O','ms','Khtml'], el = $("<p/>")[0];

   for (var i = 0; i < vendors.length; i++) {
      if (el.style[vendors[i] + 'AnimationName'] !== 'undefined') return true;
   }
})();


flowplayer(function(api, root) {

   var conf = api.conf,
      hovertimer;

   root.addClass("flowplayer").append('\
      <div class="ratio"/>\
      <div class="ui">\
         <div class="waiting"><em/><em/><em/></div>\
         <a class="fullscreen"/>\
         <a class="unload"/>\
         <div class="controls">\
            <div class="timeline">\
               <div class="buffer"/>\
               <div class="progress"/>\
            </div>\
            <div class="volume">\
               <a class="mute"></a>\
               <div class="volumeslider">\
                  <div class="volumelevel"/>\
               </div>\
            </div>\
         </div>\
         <div class="time">\
            <em class="elapsed">00:00</em>\
            <em class="remaining"/>\
            <em class="duration">99:99</em>\
         </div>\
         <div class="message"><h2/><p/></div>\
      </div>'.replace(/class="/g, 'class="fp-')
   );

   function find(klass) {
      return $(".fp-" + klass, root);
   }

   // widgets
   var progress = find("progress"),
      buffer = find("buffer"),
      elapsed = find("elapsed"),
      remaining = find("remaining"),
      waiting = find("waiting"),
      ratio = find("ratio"),
      origRatio = ratio.css("paddingTop"),

      // sliders
      timeline = find("timeline").slider(),
      timelineApi = timeline.data("api"),

      volume = find("volume"),
      volumeSlider = find("volumeslider").slider(),
      volumeApi = volumeSlider.data("api"),
      noToggle = root.hasClass("no-toggle");

   // aspect ratio
   function setRatio(val) {
      if (!parseInt(origRatio, 10))
         ratio.css("paddingTop", val * 100 + "%");

      // no inline-block support. sorry; no feature detection
      if ($.browser.msie && $.browser.version < 8) {
         $("object", root).height(root.height());
      }

   }

   function hover(flag) {
      root.toggleClass("is-mouseover", flag).toggleClass("is-mouseout", !flag);
   }

   // loading...
   if (!$.support.animation) waiting.html("<p>loading &hellip;</p>");

   setRatio(conf.ratio);

   if (noToggle) root.addClass("is-mouseover");

   // no fullscreen in IFRAME
   if (location.href != top.location.href) find("fullscreen").remove();


   api.bind("load", function() {
      root.addClass("is-loading");

   }).bind("ready", function() {

      setRatio(api.video.height / api.video.width);

      // initial time & volume
      find("duration").add(remaining).html(format(api.video.duration));
      volumeApi.slide(api.volumeLevel);


   }).bind("unload", function() {
      if (!origRatio) ratio.css("paddingTop", "");

   // buffer
   }).bind("buffer", function() {
      var video = api.video,
         max = video.buffer / video.duration;

      if (!video.seekable) timelineApi.max(max);

      buffer.animate({ width: (max * 100) + "%"}, 250, "linear");

   }).bind("buffered", function() {
      buffer.css({ width: '100%' });
      timelineApi.max(1);

   // progress
   }).bind("progress", function() {

      var time = api.video.time,
         duration = api.video.duration;

      if (!timelineApi.dragging) {
         timelineApi.slide(time / duration, api.seeking ? 0 : 250);
      }

      elapsed.html(format(time));
      remaining.html("-" + format(duration - time));

   }).bind("finish", function() {
      elapsed.html(format(api.video.duration));
      window.foo = timelineApi;
      timelineApi.slide(1, 100);

   // misc
   }).bind("beforeseek", function() {
      progress.stop();


   }).bind("volume", function() {
      volumeApi.slide(api.volumeLevel);


   }).bind("disable", function() {
      var flag = api.disabled;
      timelineApi.disable(flag);
      volumeApi.disable(flag);
      root.tc("is-disabled", api.disabled);

   }).bind("mute", function() {
      root.tc("is-muted", api.muted);

   }).bind("error", function(e, api, error) {
      root.removeClass("is-loading").addClass("is-error");

      if (error) {
         error.message = conf.errors[error.code];
         api.error = true;

         var el = $(".fp-message", root);
         $("h2", el).text(error.message);
         $("p", el).text(error.url || api.video.url || api.video.src);
      }


   // hover
   }).bind("mouseenter mouseleave", function(e) {
      if (noToggle) return;

      var is_over = e.type == "mouseenter",
         lastMove;

      // is-mouseover/out
      hover(is_over);

      if (is_over) {

         root.bind("pause.x mousemove.x volume.x", function() {
            hover(true);
            lastMove = new Date;
         });

         hovertimer = setInterval(function() {
            if (new Date - lastMove > 5000) {
               hover(false)
               lastMove = new Date;
            }
         }, 100);

      } else {
         root.unbind(".x");
         clearInterval(hovertimer);
      }


   // allow dragging over the player edge
   }).bind("mouseleave", function() {
      if (timelineApi.dragging || volumeApi.dragging) root.addClass("is-mouseover");

   // click
   }).bind("click.player", function(e) {
      if ($(e.target).is(".fp-ui, .fp-engine") || e.flash) {
         e.preventDefault();
         return api.toggle();
      }
   });

   /* controlbar elements */
   $.each(['mute', 'fullscreen', 'unload'], function(i, key) {
      find(key).click(function() {
         api[key]();
      });
   });

   timeline.bind("slide", function(e, val) {
      api.seeking = true;
      api.seek(val * api.video.duration);
   });

   volumeSlider.bind("slide", function(e, val) {
      api.volume(val);
   });

   // times
   find("time").click(function(e) {
      $(this).tc("is-inverted");
   });

   hover(false);

});

var focused;

flowplayer(function(api, root) {

   // no keyboard configured
   if (!api.conf.keyboard) return;

   // keyboard. single global listener
   $(document).bind("keydown.fp", function(e) {

      var el = focused,
         metaKeyPressed = e.ctrlKey || e.metaKey || e.altKey || e.shiftKey;

      if (el && !metaKeyPressed && api.video.time) {

         var delta = el.video.duration * 0.1;

         switch (e.which) {
            case 38: case 75: el.volume(el.volumeLevel + 0.15); break;        // volume up
            case 40: case 74: el.volume(el.volumeLevel - 0.15); break;        // volume down
            case 39: case 76: el.seeking = true; el.seek(el.video.time + delta); break;          // forward
            case 37: case 72: el.seeking = true; el.seek(el.video.time - delta); break;          // backward
            case 32: el.toggle(); break;                                      // spacebar
            case 70: el.fullscreen(); break;                            // fullscreenToggle
            case 77: el.mute(); break;                                        // mute
            case 27: el[el.isFullscreen ? 'fullscreen' : 'unload'](); break;  // esc
         }

         e.preventDefault();
      }

   });

   // hover
   root.bind("mouseenter mouseleave", function(e) {
      focused = !api.disabled && e.type == 'mouseenter' ? api : 0;
   });

});

var VENDOR = $.browser.mozilla ? "moz": "webkit",
   FS_ENTER = "fullscreen",
   FS_EXIT = "fullscreen-exit",
   FULL_PLAYER;

// detect native fullscreen support
$.support.fullscreen = typeof document.webkitCancelFullScreen == 'function' || document.mozFullScreenEnabled;


// esc button
$(document).bind(VENDOR + "fullscreenchange", function(e) {
   var el = $(document.webkitCurrentFullScreenElement || document.mozFullScreenElement);

   if (el.length) {
      FULL_PLAYER = el.trigger(FS_ENTER, [el]);
   } else {
      FULL_PLAYER.trigger(FS_EXIT, [FULL_PLAYER]);
   }

});


flowplayer(function(player, root) {

   player.isFullscreen = false;

   player.fullscreen = function(flag) {

      if (flag === undefined) flag = !player.isFullscreen;

      if ($.support.fullscreen) {

         if (flag) {
            root[0][VENDOR + 'RequestFullScreen']();
         } else {
            document[VENDOR + 'CancelFullScreen']();
         }

      } else {
         player.trigger(flag ? FS_ENTER : FS_EXIT, [player])
      }

      return player;
   };

   player.bind(FS_ENTER, function(e) {
      root.addClass("is-fullscreen");
      player.isFullscreen = true;

   }).bind(FS_EXIT, function(e) {
      root.removeClass("is-fullscreen");
      player.isFullscreen = false;
   });

   var origH = root.height();

   // handle Flash object aspect ratio on fullscreen
   player.bind("fullscreen", function() {

      var screenW = $.support.fullscreen ? screen.width : $(window).width(),
         screenH = $.support.fullscreen ? screen.height : $(window).height(),
         h = player.video.height / player.video.width * screenW;

      $("object", root).css({ height: h, marginTop: (screenH - h - 20) / 2 });

   }).bind("fullscreen-exit", function() {
      var ie7 = $.browser.msie && $.browser.version < 8;
      $("object", root).css({ height: ie7 ? origH : '', marginTop: '' })
   });

});


flowplayer(function(player, root) {

   var conf = $.extend({ active: 'is-active', advance: true, query: ".fp-playlist a" }, player.conf),
      klass = conf.active;

   // getters
   function els() {
      return $(conf.query, root);
   }

   function active() {
      return $(conf.query + "." + klass, root);
   }

   // click -> play
   var items = els().live("click", function(e) {
      var el = $(this);
      el.is("." + klass) ? player.toggle() : player.load(el.attr("href"));
      e.preventDefault();
   });

   if (items.length) {

      // highlight
      player.bind("load", function() {

         // active
         var prev = active().removeClass(klass),
            el = $("a[href*='" + player.video.src.replace(TYPE_RE, "") + "']", root).addClass(klass);

         // index
         root.removeClass("video" + els().index(prev)).addClass("video" + els().index(el));

      // without namespace callback called only once. unknown rason.
      }).bind("unload.pl", function() {
         active().tc(klass);

      });

      // next / prev
      $.each(['next', 'prev'], function(i, key) {

         player[key] = function(e) {
            var el = active()[key](); // active().next();
            if (!el.length) el = els().filter(key == 'next' ? ':first' : ':last');;
            el.click();
            e && e.preventDefault();
         };

         $(".fp-" + key, root).click(player[key]);
      });

      if (conf.advance) {
         root.unbind("finish.pl").bind("finish.pl", function() {
            root.addClass("is-playing"); // hide play button
            player.next();
         });
      }

   }


});

var CUE_RE = / ?cue\d+ ?/;

flowplayer(function(player, root) {

   var cues = player.conf.cuepoints;

   // no cuepoints
   if (!cues) return;

   function setClass(index) {
      root[0].className = root[0].className.replace(CUE_RE, " ");
      if (index >= 0) root.addClass("cue" + index);
   }

   player.bind("progress", function() {

      for (var i = 0, cue, time = player.video.time; i < cues.length; i++) {

         cue = cues[i];
         if (1 * cue) cue = { time: cue }
         if (cue.time < 0) cue.time = player.video.duration + cue.time;
         cue.index = i;

         // progress_interval / 2 = 0.125
         if (Math.abs(cue.time - time) < 0.125) {
            setClass(i);
            root.trigger("cuepoint", [player, cue]);
         }

      }

   // no CSS class name
   }).bind("unload seek", setClass);

   if (player.conf.generate_cuepoints) {

      player.bind("ready", function() {

         var duration = player.video.duration,
            timeline = $(".fp-timeline", root);

         $.each(cues, function(i, cue) {

            var time = cue.time || cue;
            if (time < 0) time = duration + cue;

            var el = $("<a/>").addClass("fp-cuepoint fp-cuepoint" + i)
               .css("left", (time / duration * 100) + "%");

            el.appendTo(timeline).mousedown(function() {
               player.seek(time);

               // preventDefault() doesn't work
               return false;
            });

         });

      });

   }

});
flowplayer(function(player, root) {

   var id = player.conf.analytics, time = 0, last = 0;

   if (id && typeof _gat !== 'undefined') {

      function track(e) {

         if (time) {
            var tracker = _gat._getTracker(id),
               video = player.video;

            tracker._setAllowLinker(true);

            // http://code.google.com/apis/analytics/docs/tracking/eventTrackerGuide.html
            tracker._trackEvent(
               "Video / Seconds played",
               player.engine + "/" + video.type,
               root.attr("title") || video.src.split("/").slice(-1)[0].replace(TYPE_RE, ''),
               Math.round(time / 1000)
            );

            time = 0;

         }

      }

      player.bind("load unload", track).bind("progress", function() {

         if (!player.seeking) {
            time += last ? (+new Date - last) : 0;
            last = +new Date;
         }

      }).bind("pause", function() {
         last = 0;
      });

      $(window).unload(track);

   }

});
/*
   Bunch of hacks to gain mobile WebKit support. Main shortomings include:

   1. cannot insert video tag dynamically -> splash screen is tricky / hack
   2. autoplay not supported

   Both of these issues cannot be feature detected. More issues can be found here:

   http://blog.millermedeiros.com/2011/03/html5-video-issues-on-the-ipad-and-how-to-solve-them/
*/

if (/iPad|MeeGo/.test(navigator.userAgent)) {

   //Warning: This is a hack!. iPad is the new IE for developers.

   flowplayer(function(player, root) {

      // custom loaded event
      var loaded;

      player.splash = player.conf.splash = false;
      player.conf.autoplay = false;

      root.bind("load", function() {
         root.addClass("is-ipad is-paused").removeClass("is-loading");
         player.ready = player.paused = true;
         player.loading = false;
         player.resume();

         // fake ready event on start
         $("video", root)[0].addEventListener("canplay", function(e) {
            root.trigger("ready").trigger("resume");

         }, false);

      });

      // force playback start with a first click
      root.bind("touchstart", function(e) {
         if (!loaded) {
            root.triggerHandler("click.player");
            loaded = true;
         }

         // fake mouseover effect with click
         if (player.playing && !root.hasClass("is-mouseover")) {
            root.addClass("is-mouseover");
            return false;
         }

      });


      player.unload = function() {
         player.pause();
         root.trigger("unload");
         loaded = false;
      };

   });

}

flowplayer(function(player, root) {

   // no embedding
   if (player.conf.embed === false) return;

   var conf = player.conf,
      ui = $(".fp-ui", root),
      trigger = $("<a/>", { "class": "fp-embed", title: 'Copy to your site'}).appendTo(ui),
      target = $("<div/>", { 'class': 'fp-embed-code'})
         .append("<label>Paste this HTML code on your site to embed.</label><textarea/>").appendTo(ui),
      area = $("textarea", target);

   player.embedCode = function() {

      var video = player.video,
         width = video.width || root.width(),
         height = video.height || root.height(),
         el = $("<div/>", { 'class': 'flowplayer', css: { width: width, height: height }}),
         tag = $("<video/>").appendTo(el);

      // configuration
      $.each(['origin', 'analytics', 'logo', 'key'], function(i, key) {
         if (conf[key]) el.attr("data-" + key, conf[key]);
      });

      // sources
      $.each(video.sources, function(i, src) {
         tag.append($("<source/>", { type: "video/" + src.type, src: src.src }));
      });

      var code = $("<foo/>", { src: "http://embed.flowplayer.org/5.0.0/embed.min.js" }).append(el);
      return $("<p/>").append(code).html().replace(/<(\/?)foo/g, "<$1script");
   };

   root.fptip(".fp-embed", "is-embedding");

   area.click(function() {
      this.select();
   });

   trigger.click(function() {
      area.text(player.embedCode());
      area[0].focus();
      area[0].select();
   });

});


$.fn.fptip = function(trigger, active) {

   return this.each(function() {

      var root = $(this);

      function close() {
         root.removeClass(active);
         $(document).unbind(".st");
      }

      $(trigger || "a", this).click(function(e) {

         e.preventDefault();

         root.tc(active);

         if (root.hasClass(active)) {

            $(document).bind("keydown.st", function(e) {
               if (e.which == 27) close();

            // click:close
            }).bind("click.st", function(e) {
               if (!$(e.target).parents("." + active).length) close();
            });
         }

      });

   });

};

}();
flowplayer(function(a,b){function i(a){var b=$("<a/>")[0];return b.href=a,b.hostname}var c=a.conf,d=c.swf.indexOf("flowplayer.org")&&c.e&&b.data("origin"),e=d?i(d):location.hostname,f=c.key;location.protocol=="file:"&&(e="localhost"),a.load.ed=1,c.hostname=e,c.origin=d||location.href,d&&b.addClass("is-embedded");if(f&&typeof key_check=="function"&&key_check(""+f,e))c.logo&&b.append($("<a>",{"class":"fp-logo",href:d}).append($("<img/>",{src:c.logo})));else{var g=$("<a/>").attr("href","http://flowplayer.org").appendTo(b),h=$(".fp-controls",b);b.bind("mouseenter mouseleave",function(b){a.ready&&g.toggle(b.type=="mouseenter")}),a.bind("progress unload",function(c){c.type=="progress"&&a.video.time<8&&a.engine!="flash"&&b.hasClass("is-mouseover")?(g.show().css({position:"absolute",left:6,bottom:h.height()+12,zIndex:99999,width:100,height:20,cursor:"pointer",backgroundImage:"url(http://stream.flowplayer.org/logo.png)"}),a.load.ed=g.is(":visible")):g.hide()})}});