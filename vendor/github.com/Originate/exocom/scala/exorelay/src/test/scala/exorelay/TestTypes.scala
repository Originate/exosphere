package exorelay

import akka.actor.ActorSystem
import akka.testkit.{ImplicitSender, TestKit}
import org.scalatest.{BeforeAndAfterAll, FlatSpecLike, Matchers}

trait CommonTest
  extends FlatSpecLike
  with Matchers

abstract class CommonAkkaTest(name: String)
  extends TestKit(ActorSystem(name))
    with CommonTest
    with BeforeAndAfterAll