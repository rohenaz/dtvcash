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
var module$contents$controllers$app_instance = null, module$exports$controllers$app = function $module$exports$controllers$app$($element$$) {
  var $$jscomp$super$this$$;
  $element$$ = void 0 === $element$$ ? document.createElement("div") : $element$$;
  null === module$contents$controllers$app_instance && (module$contents$controllers$app_instance = $$jscomp$super$this$$ = Silica.Controllers.Base.call(this, $element$$) || this);
  module$contents$controllers$app_instance.el !== $element$$ && (module$contents$controllers$app_instance.el = $element$$);
  return module$contents$controllers$app_instance;
};
$jscomp.inherits(module$exports$controllers$app, Silica.Controllers.Base);
module$exports$controllers$app.prototype.autoLoadVisible = function $module$exports$controllers$app$$autoLoadVisible$() {
  return null === document.getElementById("autoload").offsetParent;
};
module$exports$controllers$app.watchers = {autoloadVisible:function $module$exports$controllers$app$watchers$autoloadVisible$($newVal$$, $oldVal$$) {
  $newVal$$ && console.log("autoload visible!");
}};
var module$exports$Dtv = {}, module$contents$Dtv_D = {AppCntrl:module$exports$controllers$app};
Silica.compile(document);

}.call(window, window));
//# sourceMappingURL=/js/app.js.map