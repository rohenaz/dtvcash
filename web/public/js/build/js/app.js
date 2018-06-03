(function(window){
"use strict";
var $jscomp = $jscomp || {};
$jscomp.scope = {};
$jscomp.ASSUME_ES5 = !1;
$jscomp.ASSUME_NO_NATIVE_MAP = !1;
$jscomp.ASSUME_NO_NATIVE_SET = !1;
$jscomp.objectCreate = $jscomp.ASSUME_ES5 || "function" == typeof Object.create ? Object.create : function($prototype$$) {
  var $ctor$$ = function $$ctor$$$() {
  };
  $ctor$$.prototype = $prototype$$;
  return new $ctor$$;
};
$jscomp.underscoreProtoCanBeSet = function $$jscomp$underscoreProtoCanBeSet$() {
  var $x$$ = {a:!0}, $y$$ = {};
  try {
    return $y$$.__proto__ = $x$$, $y$$.a;
  } catch ($e$$) {
  }
  return !1;
};
$jscomp.setPrototypeOf = "function" == typeof Object.setPrototypeOf ? Object.setPrototypeOf : $jscomp.underscoreProtoCanBeSet() ? function($target$$, $proto$$) {
  $target$$.__proto__ = $proto$$;
  if ($target$$.__proto__ !== $proto$$) {
    throw new TypeError($target$$ + " is not extensible");
  }
  return $target$$;
} : null;
$jscomp.inherits = function $$jscomp$inherits$($childCtor$$, $parentCtor$$) {
  $childCtor$$.prototype = $jscomp.objectCreate($parentCtor$$.prototype);
  $childCtor$$.prototype.constructor = $childCtor$$;
  if ($jscomp.setPrototypeOf) {
    var $p_setPrototypeOf$$ = $jscomp.setPrototypeOf;
    $p_setPrototypeOf$$($childCtor$$, $parentCtor$$);
  } else {
    for ($p_setPrototypeOf$$ in $parentCtor$$) {
      if ("prototype" != $p_setPrototypeOf$$) {
        if (Object.defineProperties) {
          var $descriptor$$ = Object.getOwnPropertyDescriptor($parentCtor$$, $p_setPrototypeOf$$);
          $descriptor$$ && Object.defineProperty($childCtor$$, $p_setPrototypeOf$$, $descriptor$$);
        } else {
          $childCtor$$[$p_setPrototypeOf$$] = $parentCtor$$[$p_setPrototypeOf$$];
        }
      }
    }
  }
  $childCtor$$.superClass_ = $parentCtor$$.prototype;
};
$jscomp.findInternal = function $$jscomp$findInternal$($array$$, $callback$$, $thisArg$$) {
  $array$$ instanceof String && ($array$$ = String($array$$));
  for (var $len$$ = $array$$.length, $i$$ = 0; $i$$ < $len$$; $i$$++) {
    var $value$$ = $array$$[$i$$];
    if ($callback$$.call($thisArg$$, $value$$, $i$$, $array$$)) {
      return {i:$i$$, v:$value$$};
    }
  }
  return {i:-1, v:void 0};
};
$jscomp.defineProperty = $jscomp.ASSUME_ES5 || "function" == typeof Object.defineProperties ? Object.defineProperty : function($target$$, $property$$, $descriptor$$) {
  $target$$ != Array.prototype && $target$$ != Object.prototype && ($target$$[$property$$] = $descriptor$$.value);
};
$jscomp.getGlobal = function $$jscomp$getGlobal$($maybeGlobal$$) {
  return "undefined" != typeof window && window === $maybeGlobal$$ ? $maybeGlobal$$ : "undefined" != typeof global && null != global ? global : $maybeGlobal$$;
};
$jscomp.global = $jscomp.getGlobal(this);
$jscomp.polyfill = function $$jscomp$polyfill$($property$jscomp$5_split_target$$, $impl_polyfill$$, $fromLang_obj$$, $i$$) {
  if ($impl_polyfill$$) {
    $fromLang_obj$$ = $jscomp.global;
    $property$jscomp$5_split_target$$ = $property$jscomp$5_split_target$$.split(".");
    for ($i$$ = 0; $i$$ < $property$jscomp$5_split_target$$.length - 1; $i$$++) {
      var $key$$ = $property$jscomp$5_split_target$$[$i$$];
      $key$$ in $fromLang_obj$$ || ($fromLang_obj$$[$key$$] = {});
      $fromLang_obj$$ = $fromLang_obj$$[$key$$];
    }
    $property$jscomp$5_split_target$$ = $property$jscomp$5_split_target$$[$property$jscomp$5_split_target$$.length - 1];
    $i$$ = $fromLang_obj$$[$property$jscomp$5_split_target$$];
    $impl_polyfill$$ = $impl_polyfill$$($i$$);
    $impl_polyfill$$ != $i$$ && null != $impl_polyfill$$ && $jscomp.defineProperty($fromLang_obj$$, $property$jscomp$5_split_target$$, {configurable:!0, writable:!0, value:$impl_polyfill$$});
  }
};
$jscomp.polyfill("Array.prototype.find", function($orig$$) {
  return $orig$$ ? $orig$$ : function($callback$$, $opt_thisArg$$) {
    return $jscomp.findInternal(this, $callback$$, $opt_thisArg$$).v;
  };
}, "es6", "es3");
$jscomp.checkStringArgs = function $$jscomp$checkStringArgs$($thisArg$$, $arg$$, $func$$) {
  if (null == $thisArg$$) {
    throw new TypeError("The 'this' value for String.prototype." + $func$$ + " must not be null or undefined");
  }
  if ($arg$$ instanceof RegExp) {
    throw new TypeError("First argument to String.prototype." + $func$$ + " must not be a regular expression");
  }
  return $thisArg$$ + "";
};
$jscomp.polyfill("String.prototype.endsWith", function($orig$$) {
  return $orig$$ ? $orig$$ : function($searchString$$, $i$jscomp$5_opt_position$$) {
    var $string$$ = $jscomp.checkStringArgs(this, $searchString$$, "endsWith");
    $searchString$$ += "";
    void 0 === $i$jscomp$5_opt_position$$ && ($i$jscomp$5_opt_position$$ = $string$$.length);
    $i$jscomp$5_opt_position$$ = Math.max(0, Math.min($i$jscomp$5_opt_position$$ | 0, $string$$.length));
    for (var $j$$ = $searchString$$.length; 0 < $j$$ && 0 < $i$jscomp$5_opt_position$$;) {
      if ($string$$[--$i$jscomp$5_opt_position$$] != $searchString$$[--$j$$]) {
        return !1;
      }
    }
    return 0 >= $j$$;
  };
}, "es6", "es3");
var module$contents$controllers$app_instance = null, module$exports$controllers$app = function $module$exports$controllers$app$($element$$) {
  var $$jscomp$super$this$$;
  $element$$ = void 0 === $element$$ ? document.createElement("div") : $element$$;
  null === module$contents$controllers$app_instance && (module$contents$controllers$app_instance = $$jscomp$super$this$$ = Silica.Controllers.Base.call(this, $element$$) || this, $$jscomp$super$this$$.loadMore = !1, $$jscomp$super$this$$.loadingMore = !1, window.onscroll = function $window$onscroll$() {
    $$jscomp$this$$.autoLoadVisible() !== $$jscomp$this$$.loadMore && Silica.apply(function() {
      $$jscomp$this$$.loadMore = !0;
    });
  });
  var $$jscomp$this$$ = $$jscomp$super$this$$;
  module$contents$controllers$app_instance.el !== $element$$ && (module$contents$controllers$app_instance.el = $element$$);
  return module$contents$controllers$app_instance;
};
$jscomp.inherits(module$exports$controllers$app, Silica.Controllers.Base);
module$exports$controllers$app.prototype.autoLoadVisible = function $module$exports$controllers$app$$autoLoadVisible$() {
  var $elem$$ = document.getElementById("autoload");
  return $elem$$ ? this.isElementInViewport($elem$$) : !1;
};
module$exports$controllers$app.prototype.isElementInViewport = function $module$exports$controllers$app$$isElementInViewport$($el_rect$$) {
  $el_rect$$ = $el_rect$$.getBoundingClientRect();
  return 0 <= $el_rect$$.top && 0 <= $el_rect$$.left && $el_rect$$.bottom <= (window.innerHeight || document.documentElement.clientHeight) && $el_rect$$.right <= (window.innerWidth || document.documentElement.clientWidth);
};
module$exports$controllers$app.prototype.startWebtorrent = function $module$exports$controllers$app$$startWebtorrent$($client_el$$, $id$$, $magnet$$) {
  $client_el$$ = [];
  $client_el$$[$id$$] = new WebTorrent;
  $client_el$$[$id$$].add($magnet$$ + "&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com", function($file_torrent$$) {
    ($file_torrent$$ = $file_torrent$$.files.find(function($file$$) {
      var $mp4$$ = $file$$.name.endsWith(".mp4");
      $file$$ = $file$$.name.endsWith(".mp3");
      return $mp4$$ || $file$$;
    })) ? $file_torrent$$.appendTo("#message-" + $id$$) : document.getElementById("message-" + $id$$).innerHTML = "No mp4 found in this torrent";
  });
};
module$exports$controllers$app.watchers = {loadMore:function $module$exports$controllers$app$watchers$loadMore$($newOffset_newVal_offset$$, $oldVal$$) {
  $newOffset_newVal_offset$$ && (this.LoadingMore = !0, console.log("autoload visible!"), $newOffset_newVal_offset$$ = window.location.search.split("offset=")[1] || 0, $newOffset_newVal_offset$$ = parseInt($newOffset_newVal_offset$$) + 25, window.location.href = "http://dtv.cash/feed?offset=" + $newOffset_newVal_offset$$);
}};
var module$exports$Dtv = {}, module$contents$Dtv_Dtv = {AppCntrl:module$exports$controllers$app};
window.Dtv = module$contents$Dtv_Dtv;
Silica.setContext("Dtv");
Silica.compile(document);
Silica.apply(function() {
  console.log("must call apply");
});

}.call(window, window));
//# sourceMappingURL=/js/app.js.map