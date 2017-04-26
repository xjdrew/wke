package wke

/*
#include "wke.h"

extern wkeJSValue jsNativeCall(char* name, wkeJSState* es);

wkeJSValue JS_CALL gogate(wkeJSState* es) {
    const utf8* name = wkeJSToTempString(es, wkeJSParam(es,0));
    return jsNativeCall((char*)name, es);
}
*/
import "C"
