<table>
  <tr>
    <td><a href="03_communication.md">&lt;&lt; communication</a></td>
    <th>Message-Oriented Programming</th>
    <td><a href="05_todo_service.md">building the Todo service&gt;&gt;</a></td>
  </tr>
</table>


# Type checking messages

<table>
  <tr>
    <td>
      <b><i>
        Status: alpha - general idea implemented, needs feedback
      </i></b>
    </td>
  </tr>
</table>

In [chapter 2-6](06_communication.md),
we mentioned that messages are calls across functional boundaries
within your application.

To support this as much as possible,
they are built as statically typed, shared, immutable data structures.

