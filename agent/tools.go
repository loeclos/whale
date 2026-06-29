package agent

type ToolFuncResponse struct {
	Message string
	Result  any
}

type ToolChunk struct {
	ID              string
	ToolName        string
	ToolParams      string
	ApproveToolChan chan bool
}

type ToolFailedChunk struct {
	ID         string
	ToolName   string
	ToolParams string
	Reason     string
}

type ToolFunc func(...any) (ToolFuncResponse, error)

type Tools map[string]ToolFunc

// ToolManager is the function that... well, manages tool
// calls. It recieves the tool chunk (for the tool information),
// the tool function (to be able to run the tool), and the
// chunk channel to be able to communicate with the main
// go routine.
func ToolManager(tool ToolChunk, toolFunc ToolFunc, chunkChan chan any) {

	// Here we use the ApproveToolChan located inside the ToolChunk
	// to wait for the the user to confirm (or reject) the tool.

	select {
	case approved := <-tool.ApproveToolChan:
		if approved {

			// If the user approves the tool, then we run the tool
			// function.

			response, err := toolFunc(tool.ToolParams)

			// If there is an error, send it to the client through
			// the chunk chan.

			if err != nil {
				chunkChan <- ToolFailedChunk{
					ID:         tool.ID,
					ToolName:   tool.ToolName,
					ToolParams: tool.ToolParams,
					Reason:     err.Error(),
				}

				return
			}

			// If there is no error, then we send the tool response
			// back to the client, where is can be further managed
			// (such as add to chat history for the model to see).

			chunkChan <- response

			return
		} else {

			// If the user rejected the tool, then we simply send
			// a ToolFailedChunk with Reason of "tool rejected".

			chunkChan <- ToolFailedChunk{
				ID:         tool.ID,
				ToolName:   tool.ToolName,
				ToolParams: tool.ToolParams,
				Reason:     "tool rejected",
			}

			return
		}
	}
}
