package agent

// This is the type for the tool response.
// Here we have a string `Message` for relaying
// success or error messages, and Result,
// which can be of any type; this is so the function
// can return a struct of anytype.

type ToolResponse struct {
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

type Tools map[string]func(...any) (ToolResponse, error)

func ToolManager(tool ToolChunk, tools Tools, chunkChan chan any) {
	select {
	case approved := <-tool.ApproveToolChan:
		if approved {
			response, err := tools[tool.ToolName](tool.ToolParams)

			if err != nil {
				chunkChan <- ToolFailedChunk{
					ID:         tool.ID,
					ToolName:   tool.ToolName,
					ToolParams: tool.ToolParams,
					Reason:     err.Error(),
				}

				return
			}

			chunkChan <- response

			return
		} else {
			chunkChan <- ToolFailedChunk{
				ID:         tool.ID,
				ToolName:   tool.ToolName,
				ToolParams: tool.ToolParams,
				Reason:     "tool rejected",
			}
		}
	}
}
