package wke

/*
#include "wke.h"

extern wkeJSValue goNativeCall(char* name, wkeJSState* es);
extern void goTitleChanged(wkeWebView*, const char*);
extern void goURLChanged(wkeWebView*, const char*);

wkeJSValue JS_CALL gogate(wkeJSState* es) {
    const utf8* name = wkeJSToTempString(es, wkeJSParam(es,0));
    return goNativeCall((char*)name, es);
}

void titleChangedCallback(wkeWebView* webView, void* param, const wkeString* title) {
    goTitleChanged(webView, (const char*)wkeGetString(title));
}

void urlChangedCallback(wkeWebView* webView, void* param, const wkeString* url) {
    goURLChanged(webView, (const char*)wkeGetString(url));
}
*/
import "C"
