package log

import "github.com/juju/loggo"

// Logger is the global logger object
var Logger = loggo.GetLogger("goemu")

// Errorf logs a printf-formatted message at Error level
var Errorf = Logger.Errorf

// Warningf logs a printf-formatted message at Warning level
var Warningf = Logger.Warningf

// Debugf logs a printf-formatted message at Debug level
var Debugf = Logger.Debugf

// Tracef logs a printf-formatted message at Trace level
var Tracef = Logger.Tracef
