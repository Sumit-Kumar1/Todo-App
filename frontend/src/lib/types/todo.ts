export type Task = {
	Id: string;
	Title: string;
	Description: string;
	DueDate: string;
	IsDone: boolean;
};

export type TaskResponse = {
	id: string;
	title: string;
	description: string;
	dueDate: string;
	isDone: boolean;
	addedAt: string;
	modifiedAt: string;
};

export type CreateTaskRequest = {
	title: string;
	description: string;
	dueDate: string;
};
