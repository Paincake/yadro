package event

type ClientPresentError struct{}

func (e ClientPresentError) Error() string {
	return "YouShallNotPass"
}

type ClientNotPresentError struct{}

func (e ClientNotPresentError) Error() string {
	return "ClientUnknown"
}

type ClubClosedError struct{}

func (e ClubClosedError) Error() string {
	return "NotOpenYet"
}

type NoTableAvailableError struct{}

func (e NoTableAvailableError) Error() string {
	return "PlaceIsBusy"
}

type WaitingError struct{}

func (e WaitingError) Error() string {
	return "ICanWaitNoLonger"
}
