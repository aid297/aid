export let removeSpaces = (input: string): string => input.replace(/\s+/g, '')

export let replace = (input: string, search: string, replace: string) => input.replaceAll(search, replace)
