package nachos

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
)

type Route struct {
	// Path is the path of the route.
	// Appending a period in prefix or suffix is not recommended since the translation layer will auto-append it to children.
	// Furthermore, you should follow the naming convention for Nats.io (https://docs.nats.io/nats-concepts/subjects#characters-allowed-for-subject-names).
	Path string
	// BeforeAction is all the actions that should be done before the main action is performed. You can further prevent execution of
	// the main action by returning false into the action.
	BeforeAction []BeforeActionEvent
	// Action is the main action that will be executed after BeforeAction.
	Action ActionEvent
	// EndAction is all the actions that should be done after the main action is complete. You shouldn't interact with the message
	// at this stage.
	EndAction []EndActionEvent
	// Children are additional children routes that will inherit the BeforeAction and EndAction of this parent node. Having children
	// routes will automatically ignore Action since the route will be considered as a parent node.
	Children []Route
	// QueueGroup configures which queue group this route will be listening on, you can set this as nil to subscribe normally.
	// Children nodes will also inherit this if there is no override.
	QueueGroup *QueueGroup
}

type QueueGroup struct {
	Enabled bool
	Name    string
}

// Routes is where all the routes of the application is stored. You shouldn't modify this variable, it is exposed
// to enable some leeway in debugging or manipulation if needed.
var Routes = make(map[string]Route)
var Subscriptions = make(map[string]*nats.Subscription)

// TraceRoutes forces the translation layer of subscribe to log all the routes after appending them into the
// routes for debugging purposes.
var TraceRoutes = false

// Subscribe translates or flattens one or more routes before appending them into the application. In order to make
// the core system simpler, children routes are flattened by preparing the route path and inheriting all its parents
// before and after actions before being appended into the routes.
func Subscribe(routes ...Route) {
	for _, parent := range routes {
		if parent.Children == nil {
			traceAndAppend(parent)
			log.Println("Appended ", parent.Path)
			continue
		}
		translateChildren(parent)
	}
}

func translateChildren(parent Route) {
	for _, child := range parent.Children {
		var beforeActions = child.BeforeAction
		if parent.BeforeAction != nil && child.BeforeAction != nil {
			beforeActions = make([]BeforeActionEvent, len(parent.BeforeAction)+len(child.BeforeAction)-1)
			copy(beforeActions, parent.BeforeAction)
			beforeActions = append(beforeActions, child.BeforeAction...)
		}
		if parent.BeforeAction != nil && child.BeforeAction == nil {
			beforeActions = parent.BeforeAction
		}
		child.BeforeAction = beforeActions

		var endActions = child.EndAction
		if parent.EndAction != nil && child.EndAction != nil {
			endActions = make([]EndActionEvent, len(parent.EndAction)+len(child.EndAction)-1)
			copy(endActions, parent.EndAction)
			endActions = append(endActions, child.EndAction...)
		}
		if parent.EndAction != nil && child.EndAction == nil {
			endActions = parent.EndAction
		}
		child.EndAction = endActions

		queueGroup := child.QueueGroup
		if parent.QueueGroup != nil {
			queueGroup = parent.QueueGroup
		}
		child.QueueGroup = queueGroup

		child.Path = fmt.Sprint(strings.TrimSuffix(parent.Path, "."), ".", strings.TrimPrefix(child.Path, "."))
		if child.Children == nil {
			traceAndAppend(child)
			continue
		}
		translateChildren(child)
	}
}
func traceAndAppend(route Route) {
	Routes[route.Path] = route
	if TraceRoutes {
		Logger(Trace, "Adding route ", route.Path, " with ", len(route.BeforeAction), " before actions, ", len(route.EndAction), " end actions and pointer ", &route)
	}
}
