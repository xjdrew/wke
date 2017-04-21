package wke

/*
#include "wke.h"

extern jsValue jsNativeCall(char* name, jsExecState e);

jsValue JS_CALL gogate(jsExecState es)
{
    const utf8* name = jsToString(es, jsArg(es,0));
    return jsNativeCall((char*)name, es);
}
*/
import "C"
