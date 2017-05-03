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

// file callback
extern void* goFileOpen(const char*);
extern void goFileClose(void*);
extern size_t goFileSize(void*);
extern int goFileRead(void*,void*,size_t);
extern int goFileSeek(void*, int, int);

void* file_open(const char* path) {
    return goFileOpen(path);
}

void file_close(void* handle) {
    goFileClose(handle);
}

size_t file_size(void* handle) {
    return goFileSize(handle);
}

int file_read(void* handle, void* buffer, size_t size) {
    return goFileRead(handle, buffer, size);
}

int file_seek(void* handle, int offset, int origin) {
    return goFileSeek(handle, offset, origin);
}
*/
import "C"
