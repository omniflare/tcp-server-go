package main

func (r *room) broadcast(sender *client, msg string) {
	for _, m := range r.members {
		// if addr != sender.conn.RemoteAddr() {
			m.msg(msg)
		// }
	}
}
