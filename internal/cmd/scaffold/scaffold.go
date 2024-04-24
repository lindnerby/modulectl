package scaffold

func RunScaffold(opts Options) error {
	if err := opts.validate(); err != nil {
		return err
	}

	return nil
}
