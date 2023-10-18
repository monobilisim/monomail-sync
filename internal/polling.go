package internal

func GetPollingData(index int) PageData {
	if queue.Len() == 0 {
		return PageData{}
	}

	tasks := getPageByIndex(index)

	data := PageData{
		Index: index,
		Tasks: tasks,
	}

	return data
}
